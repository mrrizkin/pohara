package inertia

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/core/session"
	"github.com/mrrizkin/pohara/modules/neoweb/vite"
	"github.com/mrrizkin/pohara/resources/inertia"
)

// Inertia wraps gonertia.Inertia with additional Vite integration
type Inertia struct {
	core    *gonertia.Inertia
	flash   *SimpleFlashProvider
	session *session.Session
}

// Dependencies defines the required dependencies for Inertia
type Dependencies struct {
	fx.In

	Session *session.Session
	Config  *config.Config
	Vite    *vite.Vite
}

// New creates a new instance of Inertia with the provided dependencies
func New(deps Dependencies) (*Inertia, error) {
	options := make([]gonertia.Option, 0)

	flash := NewSimpleFlashProvider()

	options = append(options, gonertia.WithFlashProvider(flash))

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

	r, err := inertia.Entry.Open(deps.Config.Inertia.EntryPath)
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
		core:    i,
		session: deps.Session,
		flash:   flash,
	}, nil
}

func getVersionFromManifest(path string) string {
	// Try to open and read the file
	data, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, generate random string
		return hash.NanoID()
	}

	// Create MD5 hash
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// EncryptHistory enables history encryption
func (i *Inertia) EncryptHistory(ctx *fiber.Ctx) error {
	c, ok := ctx.Locals("inertia_context").(context.Context)
	if !ok {
		r, err := adaptor.ConvertRequest(ctx, true)
		if err != nil {
			return err
		}
		c = r.Context()
	}

	c = gonertia.SetEncryptHistory(c)
	ctx.Locals("inertia_context", c)
	return nil
}

// ClearHistory clears the history
func (i *Inertia) ClearHistory(ctx *fiber.Ctx) error {
	c, ok := ctx.Locals("inertia_context").(context.Context)
	if !ok {
		r, err := adaptor.ConvertRequest(ctx, true)
		if err != nil {
			return err
		}
		c = r.Context()
	}

	c = gonertia.ClearHistory(c)
	ctx.Locals("inertia_context", c)
	return nil
}

// EncryptHistoryMiddleware provides middleware for encrypting history
func (i *Inertia) EncryptHistoryMiddleware(ctx *fiber.Ctx) error {
	if err := i.EncryptHistory(ctx); err != nil {
		return err
	}

	return ctx.Next()
}

// Middleware provides Inertia middleware for Fiber
func (i *Inertia) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var next bool

		ses, err := i.session.Get(c)
		if err != nil {
			return err
		}

		c.Locals("session_id", ses.ID())
		i.core.ShareProp("flash", fiber.Map{
			"message": ses.Get("message"),
		})

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next = true
			// Convert again in case request may modify by middleware
			c.Request().Header.SetMethod(r.Method)
			c.Request().SetRequestURI(r.RequestURI)
			c.Request().SetHost(r.Host)
			c.Request().Header.SetHost(r.Host)
			for key, val := range r.Header {
				for _, v := range val {
					c.Request().Header.Set(key, v)
				}
			}
			adaptor.CopyContextToFiberContext(r.Context(), c.Context())
		})

		if err := adaptor.HTTPHandler(i.core.Middleware(nextHandler))(c); err != nil {
			return err
		}

		if next {
			return c.Next()
		}
		return nil
	}
}

// Redirect redirects to the given URL
func (i *Inertia) Redirect(ctx *fiber.Ctx, url string, status ...int) error {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return err
	}

	if c, ok := ctx.Locals("inertia_context").(context.Context); ok {
		r = r.WithContext(c)
	}

	w := newResponseWriter()
	i.core.Redirect(w, r, url, status...)
	return writeResponse(ctx, w)
}

// Location redirect to the given external URL
func (i *Inertia) Location(ctx *fiber.Ctx, url string, status ...int) error {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return err
	}

	if c, ok := ctx.Locals("inertia_context").(context.Context); ok {
		r = r.WithContext(c)
	}

	w := newResponseWriter()
	i.core.Location(w, r, url, status...)
	return writeResponse(ctx, w)
}

// Back redirects to the previous URL
func (i *Inertia) Back(ctx *fiber.Ctx, status ...int) error {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return err
	}

	if c, ok := ctx.Locals("inertia_context").(context.Context); ok {
		r = r.WithContext(c)
	}

	w := newResponseWriter()
	i.core.Back(w, r, status...)
	return writeResponse(ctx, w)
}

// Render renders an Inertia component with the given props
func (i *Inertia) Render(ctx *fiber.Ctx, component string, props ...gonertia.Props) error {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return err
	}

	if c, ok := ctx.Locals("inertia_context").(context.Context); ok {
		r = r.WithContext(c)
	}

	c := r.Context()
	if _, ok := c.Value("session_id").(string); !ok {
		sess, err := i.session.Get(ctx)
		if err != nil {
			return err
		}

		c = context.WithValue(c, "session_id", sess.ID())
	}

	if flashError, err := i.flash.GetErrors(c); err == nil {
		c = gonertia.SetValidationErrors(c, flashError)
	}

	r = r.WithContext(c)
	w := newResponseWriter()
	if err := i.core.Render(w, r, component, props...); err != nil {
		return err
	}

	return writeResponse(ctx.Type("html"), w)
}

func (i *Inertia) IsInertiaRequest(ctx *fiber.Ctx) bool {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return false
	}

	return gonertia.IsInertiaRequest(r)
}

func (i *Inertia) AddValidationError(ctx *fiber.Ctx, errMap fiber.Map) error {
	c, ok := ctx.Locals("inertia_context").(context.Context)
	if !ok {
		r, err := adaptor.ConvertRequest(ctx, true)
		if err != nil {
			return err
		}
		c = r.Context()

	}

	if _, ok = c.Value("session_id").(string); !ok {
		sess, err := i.session.Get(ctx)
		if err != nil {
			return err
		}

		c = context.WithValue(c, "session_id", sess.ID())
	}

	c = gonertia.AddValidationErrors(c, gonertia.ValidationErrors(errMap))
	ctx.Locals("inertia_context", c)
	return nil
}
