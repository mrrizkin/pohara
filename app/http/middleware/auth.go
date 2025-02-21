package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/app/service"
	"github.com/mrrizkin/pohara/modules/inertia"
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/mrrizkin/pohara/modules/session"
	"go.uber.org/fx"
)

type AuthMiddleware struct {
	authService  *service.AuthService
	inertia      *inertia.Inertia
	log          *logger.Logger
	sessionStore *session.Store
}

type AuthMiddlewareDependencies struct {
	fx.In

	AuthService  *service.AuthService
	Inertia      *inertia.Inertia
	SessionStore *session.Store
	Logger       *logger.Logger
}

func NewAuthMiddleware(deps AuthMiddlewareDependencies) *AuthMiddleware {
	return &AuthMiddleware{
		authService:  deps.AuthService,
		inertia:      deps.Inertia,
		sessionStore: deps.SessionStore,
		log:          deps.Logger.Auth().Scope("auth_middleware"),
	}
}

func (a *AuthMiddleware) Authenticated(ctx *fiber.Ctx) error {
	if err := a.authService.Authenticated(ctx); err != nil {
		if strings.HasPrefix(ctx.Path(), "/_/") {
			return a.inertia.Render(ctx.Status(fiber.StatusUnauthorized), "error/unauthorized")
		}

		return err
	}

	return ctx.Next()
}
