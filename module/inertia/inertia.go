package inertia

import (
	"html/template"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/module/vite"
)

// Constants
const (
	defaultTemplatePath = "resources/views/root.html"
)

// Inertia wraps gonertia.Inertia with additional Vite integration
type Inertia struct {
	core *gonertia.Inertia
	vite *vite.Vite
}

// Dependencies defines the required dependencies for Inertia
type Dependencies struct {
	fx.In

	Vite   *vite.Vite
	Config *config.App
}

// Result wraps the Inertia instance for fx dependency injection
type Result struct {
	fx.Out

	Inertia *Inertia
}

// responseWriter implements http.ResponseWriter for adapting between Fiber and net/http
type responseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

// New creates a new instance of Inertia with the provided dependencies
func New(deps Dependencies) (Result, error) {
	options, err := buildOptions(deps)
	if err != nil {
		return Result{}, err
	}

	i, err := gonertia.NewFromFile(defaultTemplatePath, options...)
	if err != nil {
		return Result{}, err
	}

	registerTemplateFuncs(i, deps.Vite)

	return Result{
		Inertia: &Inertia{
			core: i,
			vite: deps.Vite,
		},
	}, nil
}

// Middleware provides Inertia middleware for Fiber
func (i *Inertia) Middleware(h fiber.Handler) fiber.Handler {
	return adaptor.HTTPHandler(i.core.Middleware(adaptor.FiberHandler(h)))
}

// Render renders an Inertia component with the given props
func (i *Inertia) Render(ctx *fiber.Ctx, component string, props ...gonertia.Props) error {
	r, err := adaptor.ConvertRequest(ctx, false)
	if err != nil {
		return err
	}

	w := newResponseWriter()
	i.core.Render(w, r, component, props...)

	return writeResponse(ctx, w)
}

// Private helper functions

func buildOptions(deps Dependencies) ([]gonertia.Option, error) {
	if !deps.Config.IsProduction() {
		return nil, nil
	}

	manifestContent, err := deps.Vite.Content()
	if err != nil {
		return nil, err
	}

	return []gonertia.Option{
		gonertia.WithVersion(string(manifestContent)),
	}, nil
}

func registerTemplateFuncs(i *gonertia.Inertia, v *vite.Vite) {
	i.ShareTemplateFunc("reactRefresh", func() template.HTML {
		return template.HTML(v.ReactRefresh())
	})
	i.ShareTemplateFunc("vite", func(input string) template.HTML {
		return template.HTML(v.Entry(input))
	})
}

func writeResponse(ctx *fiber.Ctx, w *responseWriter) error {
	for key, values := range w.Header() {
		for _, value := range values {
			ctx.Set(key, value)
		}
	}

	return ctx.Type("html").Status(w.StatusCode()).Send(w.Body())
}

// responseWriter implementation

func newResponseWriter() *responseWriter {
	return &responseWriter{
		header:     http.Header{},
		statusCode: http.StatusOK,
	}
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
	}
}

func (w *responseWriter) Body() []byte {
	return w.body
}

func (w *responseWriter) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}
