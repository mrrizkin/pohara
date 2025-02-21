package server

import (
	"embed"
	"fmt"
	"html/template"
	"strings"
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

	sysConfig "github.com/mrrizkin/pohara/modules/common/config"
	"github.com/mrrizkin/pohara/modules/common/debug"
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/mrrizkin/pohara/modules/server/cli"
	"github.com/mrrizkin/pohara/modules/server/config"
	"github.com/mrrizkin/pohara/modules/session"
)

//go:embed templates/swagger.html
var swaggerTemplates embed.FS

type Dependencies struct {
	fx.In

	Config *config.Config
	Logger *logger.Logger
}

var Module = fx.Module("server",
	sysConfig.Load(&config.Config{}),
	fx.Provide(NewServer,
		cli.NewStartServerCmd,
		cli.RegisterCommands,
	),
	fx.Decorate(SetupRouter),
)

func NewServer(deps Dependencies) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		Prefork:               deps.Config.Prefork,
		AppName:               deps.Config.Name,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return ErrorHandler(!deps.Config.IsProduction(), c, err)
		},
	})

	app.Static("/", "public")
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: deps.Logger.Request().GetLogger().(*zerolog.Logger),
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

func swagger(config *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		if config.IsProduction() {
			return c.Status(fiber.StatusNotFound).Send(nil)
		}

		tmpl, _ := template.ParseFS(swaggerTemplates, "templates/swagger.html")
		data := struct {
			SwaggerPath string
		}{
			SwaggerPath: config.SwaggerPath,
		}

		var html strings.Builder
		tmpl.Execute(&html, data)

		return c.Type("html").Send([]byte(html.String()))
	}
}
