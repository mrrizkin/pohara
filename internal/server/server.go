package server

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/fx"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog"

	"github.com/mrrizkin/pohara/config"
	debugtrace "github.com/mrrizkin/pohara/internal/common/debug-trace"
	"github.com/mrrizkin/pohara/internal/common/validator"
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/internal/session"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/mrrizkin/pohara/internal/web/template"
)

type Dependencies struct {
	fx.In

	Config *config.App
	Log    ports.Logger
}

type Result struct {
	fx.Out

	App *fiber.App
}

type Map map[string]interface{}

func New(deps Dependencies) (Result, error) {
	app := fiber.New(fiber.Config{
		Prefork:               deps.Config.APP_PREFORK,
		AppName:               deps.Config.APP_NAME,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			if c.Get("X-Requested-With") != "XMLHttpRequest" {
				if stackTrace, ok := c.Locals("stack_trace").([]debugtrace.StackFrame); ok {
					html := errorPageWithTrace(stackTrace, err, code)
					return c.Type("html").Send([]byte(html))
				}

				html := errorPage(err, code)
				return c.Type("html").Send([]byte(html))
			}

			return c.Status(code).JSON(validator.GlobalErrorResponse{
				Status: "error",
				Detail: err.Error(),
			})
		},
	})

	app.Static("/", "public")
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: deps.Log.GetLogger().(*zerolog.Logger),
	}))
	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			c.Locals("stack_trace", debugtrace.StackTrace(11))
			deps.Log.Error(fmt.Sprintf("panic: %v\n", e))
		},
	}))
	app.Use(idempotency.New())

	return Result{
		App: app,
	}, nil
}

var Module = fx.Module("server",
	fx.Provide(New),
	fx.Decorate(
		fx.Annotate(
			setupRouter,
			fx.ParamTags("", "", "", "", "", "", "", "", `group:"router"`, `group:"api_router"`),
		),
	),
	fx.Invoke(func(app *fiber.App, config *config.App, log ports.Logger) error {
		log.Info("server started", "app_name", config.APP_NAME, "port", config.APP_PORT)
		return app.Listen(fmt.Sprintf(":%d", config.APP_PORT))
	}),
)

func setupRouter(
	app *fiber.App,
	session *session.Session,
	config *config.App,
	inertia *inertia.Inertia,
	log ports.Logger,
	cache ports.Cache,
	validator *validator.Validator,
	template *template.Template,
	routes []WebRouter,
	apiRoutes []ApiRouter,
) *fiber.App {
	app.Get("/api/v1/docs/swagger", swagger(config))

	(WebRouters)(routes).Register(
		NewRouter(
			app.Group("/",
				csrf.New(csrf.Config{
					KeyLookup:         fmt.Sprintf("cookie:%s", config.CSRF_KEY),
					CookieName:        config.CSRF_COOKIE_NAME,
					CookieSameSite:    config.CSRF_SAME_SITE,
					CookieSecure:      config.CSRF_SECURE,
					CookieSessionOnly: true,
					CookieHTTPOnly:    config.CSRF_HTTP_ONLY,
					SingleUseToken:    true,
					Expiration:        time.Duration(config.CSRF_EXPIRATION) * time.Second,
					KeyGenerator:      utils.UUIDv4,
					ErrorHandler:      csrf.ConfigDefault.ErrorHandler,
					Extractor:         csrf.CsrfFromCookie(config.CSRF_KEY),
					Session:           session.Store,
					SessionKey:        "fiber.csrf.token",
					HandlerContextKey: "fiber.csrf.handler",
				}),
				cors.New(),
				helmet.New(),
			),
			template,
			inertia,
			config,
			cache,
			log,
			validator,
		))

	(ApiRouters)(apiRoutes).Register(
		NewRouter(
			app.Group("/api",
				cors.New(cors.Config{
					AllowOrigins: "*",
					AllowHeaders: "Origin, Content-Type, Accept, pohara-api-token",
				}),
			),
			template,
			inertia,
			config,
			cache,
			log,
			validator,
		))

	return app
}

func swagger(config *config.App) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		if config.IsProduction() {
			return c.Status(fiber.StatusNotFound).Send(nil)
		}

		html := fmt.Sprintf(`<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>Swagger API Reference - Scalar</title>
        <link rel="icon" type="image/svg+xml" href="https://docs.scalar.com/favicon.svg">
        <link rel="icon" type="image/png" href="https://docs.scalar.com/favicon.png">
    </head>
    <body>
        <script id="api-reference" data-url="%s"></script>
        <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    </body>
</html>`, config.SWAGGER_PATH)

		return c.Type("html").Send([]byte(html))
	}
}
