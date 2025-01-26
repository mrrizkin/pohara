package server

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/common/debug"
	"github.com/mrrizkin/pohara/modules/common/response"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/session"
)

type Dependencies struct {
	fx.In

	Config *config.Config
	Logger *logger.ZeroLog
}

var Module = fx.Module("server",
	fx.Provide(
		NewServer,
	),

	fx.Decorate(SetupRouter),
	fx.Invoke(StartServer),
)

func NewServer(deps Dependencies) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		Prefork:               deps.Config.App.Prefork,
		AppName:               deps.Config.App.Name,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			if c.Get("X-Requested-With") != "XMLHttpRequest" {
				if stackTrace, ok := c.Locals("stack_trace").([]debug.StackFrame); ok {
					html := errorPageWithTrace(stackTrace, err, code)
					return c.Type("html").Status(code).Send([]byte(html))
				}

				html := errorPage(err, code)
				return c.Type("html").Status(code).Send([]byte(html))
			}

			detail := ""
			if !deps.Config.IsProduction() {
				var stackFrames []debug.StackFrame
				if stack, ok := c.Locals("stack_trace").([]debug.StackFrame); ok {
					stackFrames = stack
				} else if stack, err := debug.StackTrace(); err == nil {
					stackFrames = stack
				}

				for _, frame := range stackFrames {
					detail += fmt.Sprintf("%s (%s:%d)\n", frame.Function, frame.File, frame.Line)
				}
			}

			return c.Status(code).JSON(response.ErrorMsg(err.Error(), detail))
		},
	})

	app.Static("/", "public")
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: deps.Logger.GetLogger().(*zerolog.Logger),
	}))
	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			if stackFrames, err := debug.StackTrace(); err == nil {
				c.Locals("stack_trace", stackFrames)
			}
			deps.Logger.Error(fmt.Sprintf("panic: %v\n", e))
		},
	}))
	app.Use(idempotency.New())

	return app, nil
}

type SetupRouterDependecies struct {
	fx.In

	App          *fiber.App
	SessionStore *session.Store
	Config       *config.Config

	WebRoutes []WebRouter `group:"web_router"`
	ApiRoutes []ApiRouter `group:"api_router"`
}

func SetupRouter(deps SetupRouterDependecies) *fiber.App {
	deps.App.Get("/api/v1/docs/swagger", swagger(deps.Config))

	(ApiRouters)(deps.ApiRoutes).Register(
		deps.App.Group("/api",
			cors.New(cors.Config{
				AllowOrigins: "*",
				AllowHeaders: "Origin, Content-Type, Accept, pohara-api-token",
			}),
		))

	(WebRouters)(deps.WebRoutes).Register(
		deps.App.Group("/",
			csrf.New(csrf.Config{
				KeyLookup:         fmt.Sprintf("cookie:%s", deps.Config.CSRF.CookieName),
				CookieName:        deps.Config.CSRF.CookieName,
				CookieSameSite:    deps.Config.CSRF.SameSite,
				CookieSecure:      deps.Config.CSRF.Secure,
				CookieSessionOnly: true,
				CookieHTTPOnly:    deps.Config.CSRF.HttpOnly,
				SingleUseToken:    true,
				Expiration:        time.Duration(deps.Config.CSRF.Expiration) * time.Second,
				KeyGenerator:      utils.UUIDv4,
				ErrorHandler:      csrf.ConfigDefault.ErrorHandler,
				Extractor:         csrf.CsrfFromCookie(deps.Config.CSRF.CookieName),
				Session:           deps.SessionStore.Store,
				SessionKey:        "fiber.csrf.token",
				HandlerContextKey: "fiber.csrf.handler",
			}),
			cors.New(),
			helmet.New(),
		))

	return deps.App
}

func StartServer(
	lx fx.Lifecycle,
	app *fiber.App,
	config *config.Config,
	log *logger.ZeroLog,
) error {
	lx.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := app.Listen(fmt.Sprintf(":%d", config.App.Port)); err != nil {
					log.Fatal("failed to start server", "error", err)
				}
			}()
			log.Info("server started", "app_name", config.App.Name, "port", config.App.Port)
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})

	return nil
}

func swagger(config *config.Config) func(c *fiber.Ctx) error {
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
</html>`, config.SwaggerPath)

		return c.Type("html").Send([]byte(html))
	}
}
