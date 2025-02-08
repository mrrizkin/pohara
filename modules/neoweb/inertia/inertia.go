package inertia

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/core/session"
	"github.com/mrrizkin/pohara/modules/neoweb/vite"
	"github.com/mrrizkin/pohara/resources"
)

type InertiaContext int

const (
	SessionID InertiaContext = iota
)

type Inertia struct {
	core         *gonertia.Inertia
	sessionStore *session.Store
}

type Dependencies struct {
	fx.In

	SessionStore *session.Store
	Config       *config.Config
	Vite         *vite.Vite
}

func New(deps Dependencies) (*Inertia, error) {
	options := []gonertia.Option{
		gonertia.WithFlashProvider(NewSimpleFlashProvider()),
	}

	if deps.Config.Inertia.ContainerID != "" {
		options = append(options, gonertia.WithContainerID(deps.Config.Inertia.ContainerID))
	}

	if deps.Config.Inertia.ManifestPath != "" {
		options = append(
			options,
			gonertia.WithVersion(getVersionFromManifest(deps.Config.Inertia.ManifestPath)),
		)
	}

	if deps.Config.Inertia.EncryptHistory {
		options = append(options, gonertia.WithEncryptHistory(deps.Config.Inertia.EncryptHistory))
	}

	r, err := resources.Admin.Open(deps.Config.Inertia.EntryPath)
	if err != nil {
		return nil, err
	}

	i, err := gonertia.NewFromReader(r, options...)
	if err != nil {
		return nil, err
	}

	i.ShareTemplateFunc("reactRefresh", func() template.HTML {
		return template.HTML(deps.Vite.ReactRefresh())
	})
	i.ShareTemplateFunc("vite", func(input string) template.HTML {
		return template.HTML(deps.Vite.Entry(input))
	})

	return &Inertia{
		core:         i,
		sessionStore: deps.SessionStore,
	}, nil
}

func getVersionFromManifest(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return hash.NanoID()
	}

	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func (i *Inertia) Middleware(c *fiber.Ctx) error {
	var next bool

	session, err := i.sessionStore.Get(c)
	if err != nil {
		return err
	}

	c.Context().SetUserValue(SessionID, session.ID())

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next = true

		c.Request().Header.SetMethod(r.Method)
		c.Request().SetRequestURI(r.RequestURI)
		c.Request().SetHost(r.Host)
		c.Request().Header.SetHost(r.Host)
		for key, val := range r.Header {
			for _, v := range val {
				c.Request().Header.Set(key, v)
			}
		}
		copyContextToFiberContext(r.Context(), c.Context())
	})

	if err := adaptor.HTTPHandler(i.core.Middleware(nextHandler))(c); err != nil {
		return err
	}
	if next {
		return c.Next()
	}

	return nil
}

func (i *Inertia) Redirect(ctx *fiber.Ctx, url string, status ...int) error {
	r, err := convertRequest(ctx)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	i.core.Redirect(w, r, url, status...)
	return writeResponse(ctx, w)
}

func (i *Inertia) Location(ctx *fiber.Ctx, url string, status ...int) error {
	r, err := convertRequest(ctx)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	i.core.Location(w, r, url, status...)
	return writeResponse(ctx, w)
}

func (i *Inertia) Back(ctx *fiber.Ctx, status ...int) error {
	r, err := convertRequest(ctx)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	i.core.Back(w, r, status...)
	return writeResponse(ctx, w)
}

func (i *Inertia) Render(ctx *fiber.Ctx, component string, props ...gonertia.Props) error {
	r, err := convertRequest(ctx)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	if err := i.core.Render(w, r, component, props...); err != nil {
		return err
	}

	return writeResponse(ctx.Type("html"), w)
}

func EncryptHistory(ctx *fiber.Ctx) error {
	ctx.SetUserContext(gonertia.SetEncryptHistory(ctx.UserContext()))
	return nil
}

func ClearHistory(ctx *fiber.Ctx) error {
	ctx.SetUserContext(gonertia.ClearHistory(ctx.UserContext()))
	return nil
}

func EncryptHistoryMiddleware(ctx *fiber.Ctx) error {
	if err := EncryptHistory(ctx); err != nil {
		return err
	}

	return ctx.Next()
}

func IsInertiaRequest(ctx *fiber.Ctx) bool {
	r, err := convertRequest(ctx)
	if err != nil {
		return false
	}

	return gonertia.IsInertiaRequest(r)
}

func AddValidationError(ctx *fiber.Ctx, errMap fiber.Map) error {
	ctx.SetUserContext(
		gonertia.AddValidationErrors(ctx.UserContext(), gonertia.ValidationErrors(errMap)),
	)
	return nil
}

func convertRequest(ctx *fiber.Ctx) (*http.Request, error) {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return nil, err
	}

	userCtx := ctx.UserContext()
	if userCtx != context.Background() {
		copyContextToFiberContext(userCtx, ctx.Context())
	}

	r = r.WithContext(ctx.Context())
	return r, nil
}

func copyContextToFiberContext(ctx interface{}, requestContext *fasthttp.RequestCtx) {
	if ctx == context.Background() {
		return
	}
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()
	if contextKeys.Kind() == reflect.Struct {
		var lastKey interface{}
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			/* #nosec */
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).
				Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "noCopy" {
				break
			} else if reflectField.Name == "Context" {
				copyContextToFiberContext(reflectValue.Interface(), requestContext)
			} else if reflectField.Name == "key" {
				lastKey = reflectValue.Interface()
			} else if lastKey != nil && reflectField.Name == "val" {
				requestContext.SetUserValue(lastKey, reflectValue.Interface())
			} else {
				lastKey = nil
			}
		}
	}
}
