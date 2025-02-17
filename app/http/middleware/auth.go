package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/inertia"
	"go.uber.org/fx"
)

type AuthMiddleware struct {
	authService *service.AuthService
	inertia     *inertia.Inertia
}

type AuthMiddlewareDependencies struct {
	fx.In

	AuthService *service.AuthService
	Inertia     *inertia.Inertia
}

func NewAuthMiddleware(deps AuthMiddlewareDependencies) *AuthMiddleware {
	return &AuthMiddleware{
		authService: deps.AuthService,
		inertia:     deps.Inertia,
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
