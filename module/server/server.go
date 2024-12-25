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
	debugtrace "github.com/mrrizkin/pohara/module/debug-trace"
	"github.com/mrrizkin/pohara/module/logger"
	"github.com/mrrizkin/pohara/module/session"
	"github.com/mrrizkin/pohara/module/validator"
	"github.com/mrrizkin/pohara/module/view"
)

type Dependencies struct {
	fx.In

	Config *config.App
	Log    *logger.Logger
}

type Result struct {
	fx.Out

	App *fiber.App
}

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
				// stackTrace := debugtrace.StackTrace()
				if stackTrace, ok := c.Locals("stack_trace").([]debugtrace.StackFrame); ok {
					errMsg := err.Error()
					markup := ""

					for _, stack := range stackTrace {
						markup = fmt.Sprintf(`%s<li><strong>%s</strong> in <em>%s</em> at line %d</li>`,
							markup, stack.Function, stack.File, stack.Line)
					}

					html := fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Error %d</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f8f9fa;
            color: #343a40;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 800px;
            margin: 50px auto;
            background: #fff;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            border-radius: 8px;
            overflow: hidden;
        }
        .header {
            background-color: #dc3545;
            color: #fff;
            padding: 20px;
            text-align: center;
        }
        .content {
            padding: 20px;
        }
        ul {
            list-style: none;
            padding: 0;
        }
        li {
            margin-bottom: 10px;
            background: #f8d7da;
            padding: 10px;
            border-left: 5px solid #dc3545;
            border-radius: 4px;
        }
        strong {
            font-weight: bold;
        }
        em {
            font-style: italic;
            color: #6c757d;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Error %d</h1>
            <p>%s</p>
        </div>
        <div class="content">
            <h2>Stack Trace:</h2>
            <ul>
                %s
            </ul>
        </div>
    </div>
</body>
</html>`, code, code, errMsg, markup)

					return c.Type("html").Send([]byte(html))
				}

				errMsg := err.Error()

				html := fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Error %d</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f8f9fa;
            color: #343a40;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 800px;
            margin: 50px auto;
            background: #fff;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            border-radius: 8px;
            overflow: hidden;
        }
        .header {
            background-color: #dc3545;
            color: #fff;
            padding: 20px;
            text-align: center;
        }
        .content {
            padding: 20px;
        }
        ul {
            list-style: none;
            padding: 0;
        }
        li {
            margin-bottom: 10px;
            background: #f8d7da;
            padding: 10px;
            border-left: 5px solid #dc3545;
            border-radius: 4px;
        }
        strong {
            font-weight: bold;
        }
        em {
            font-style: italic;
            color: #6c757d;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Error %d</h1>
            <p>%s</p>
        </div>
        
    </div>
</body>
</html>`, code, code, errMsg)

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
		fx.Annotate(setupRouter, fx.ParamTags("", "", "", `group:"router"`, `group:"api_router"`)),
	),
	fx.Invoke(func(app *fiber.App, config *config.App, log *logger.Logger, view *view.View) error {
		view.Compile()
		log.Info("server started", "app_name", config.APP_NAME, "port", config.APP_PORT)
		return app.Listen(fmt.Sprintf(":%d", config.APP_PORT))
	}),
)

func setupRouter(
	app *fiber.App,
	session *session.Session,
	config *config.App,
	routes []WebRouter,
	apiRoutes []ApiRouter,
) *fiber.App {
	app.Get("/api/v1/docs/swagger", swagger(config))

	(WebRouters)(routes).Register(
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
		))

	(ApiRouters)(apiRoutes).Register(
		app.Group("/api",
			cors.New(cors.Config{
				AllowOrigins: "*",
				AllowHeaders: "Origin, Content-Type, Accept, pohara-api-token",
			}),
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