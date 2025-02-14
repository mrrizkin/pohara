package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/core/session"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
	"go.uber.org/fx"
)

type AuthMiddleware struct {
	sessionStore *session.Store
	authService  *service.AuthService
	inertia      *inertia.Inertia
}

type AuthMiddlewareDependencies struct {
	fx.In

	SessionStore *session.Store
	AuthService  *service.AuthService
	Inertia      *inertia.Inertia
}

func NewAuthMiddleware(deps AuthMiddlewareDependencies) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore: deps.SessionStore,
		authService:  deps.AuthService,
		inertia:      deps.Inertia,
	}
}

func (a *AuthMiddleware) Authenticated(ctx *fiber.Ctx) error {
	if err := a.authService.Authenticated(ctx); err != nil {
		if errors.Is(err, fiber.ErrUnauthorized) && strings.HasPrefix(ctx.Path(), "/_/") {
			return a.inertia.Render(ctx.Status(fiber.StatusUnauthorized), "error/unauthorized")
		}

		return err
	}

	return ctx.Next()
}
