package inertia

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/web"
)

// Inertia wraps gonertia.Inertia with additional Vite integration
type Inertia struct {
	core *gonertia.Inertia
}

// Dependencies defines the required dependencies for Inertia
type Dependencies struct {
	fx.In

	Config        *config.App
	InertiaConfig *config.Inertia
}

// Result wraps the Inertia instance for fx dependency injection
type Result struct {
	fx.Out

	Inertia *Inertia
}

var Module = fx.Module("inertia",
	fx.Provide(New),
	fx.Decorate(extendInertia),
	fx.Invoke(func(*Inertia) {}),
)

// New creates a new instance of Inertia with the provided dependencies
func New(deps Dependencies) (Result, error) {
	options := make([]gonertia.Option, 0)

	if deps.InertiaConfig.CONTAINER_ID != "" {
		options = append(options, gonertia.WithContainerID(deps.InertiaConfig.CONTAINER_ID))
	}

	if deps.InertiaConfig.MANIFEST_PATH != "" {
		options = append(options, gonertia.WithVersion(getVersionFromManifest(deps.InertiaConfig.MANIFEST_PATH)))
	}

	if deps.InertiaConfig.ENCRYPT_HISTORY {
		options = append(options, gonertia.WithEncryptHistory(deps.InertiaConfig.ENCRYPT_HISTORY))
	}

	r, err := web.InertiaRoot.Open(deps.InertiaConfig.ENTRY_PATH)
	if err != nil {
		return Result{}, err
	}

	i, err := gonertia.NewFromReader(r, options...)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Inertia: &Inertia{
			core: i,
		},
	}, nil
}

func getVersionFromManifest(path string) string {
	// Try to open and read the file
	data, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, generate random string
		return randString(16)
	}

	// Create MD5 hash
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// Helper function to generate random string
func randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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

// Middleware provides Inertia middleware for Fiber
func (i *Inertia) Middleware(h fiber.Handler) fiber.Handler {
	return adaptor.HTTPHandler(i.core.Middleware(adaptor.FiberHandler(h)))
}

// Redirect redirects to the given URL
func (i *Inertia) Redirect(ctx *fiber.Ctx, url string, status int) error {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return err
	}

	if c, ok := ctx.Locals("inertia_context").(context.Context); ok {
		r.WithContext(c)
	}

	w := newResponseWriter()
	i.core.Redirect(w, r, url, status)
	return writeResponse(ctx, w)
}

// Render renders an Inertia component with the given props
func (i *Inertia) Render(ctx *fiber.Ctx, component string, props ...gonertia.Props) error {
	r, err := adaptor.ConvertRequest(ctx, true)
	if err != nil {
		return err
	}

	if c, ok := ctx.Locals("inertia_context").(context.Context); ok {
		r.WithContext(c)
	}

	w := newResponseWriter()
	if err := i.core.Render(w, r, component, props...); err != nil {
		return err
	}

	return writeResponse(ctx.Type("html"), w)
}

type ExtendResult struct {
	fx.Out

	ShareTemplateFunc map[string]interface{} `group:"inertia_template_func"`
	ShareProp         map[string]interface{} `group:"inertia_prop"`
}

func Extend(fn any) any {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("extend expects a function")
	}

	// Validate that the function returns exactly one value
	if fnType.NumOut() != 1 {
		panic("function must return 1 thing")
	}

	// Validate that the return value is a struct
	if fnType.Out(0).Kind() != reflect.Struct {
		panic("function must return a struct")
	}

	// Validate that the return type matches ExtendResult
	expectedType := reflect.TypeOf(ExtendResult{})
	if !fnType.Out(0).AssignableTo(expectedType) {
		panic("function must return ExtendResult")
	}

	return fn
}

func NewExtend(data map[string]interface{}) ExtendResult {
	templateFunc := make(map[string]interface{})
	prop := make(map[string]interface{})

	for key, value := range data {
		// Check if the value is a function using reflection
		valueType := reflect.TypeOf(value)
		if valueType != nil && valueType.Kind() == reflect.Func {
			templateFunc[key] = value
			continue
		}
		// If not a function, add to props
		prop[key] = value
	}

	return ExtendResult{
		ShareTemplateFunc: templateFunc,
		ShareProp:         prop,
	}
}

type ExtendInertiaDependencies struct {
	fx.In

	Inertia           *Inertia
	ShareTemplateFunc []map[string]interface{} `group:"inertia_template_func"`
	ShareProp         []map[string]interface{} `group:"inertia_prop"`
}

func extendInertia(deps ExtendInertiaDependencies) *Inertia {
	inertia := deps.Inertia

	for _, fns := range deps.ShareTemplateFunc {
		for name, fn := range fns {
			inertia.core.ShareTemplateFunc(name, fn)
		}
	}

	for _, props := range deps.ShareProp {
		for name, prop := range props {
			inertia.core.ShareProp(name, prop)
		}
	}

	return inertia
}
