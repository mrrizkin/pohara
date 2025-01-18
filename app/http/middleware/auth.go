package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/core/session"
	"go.uber.org/fx"
)

type AuthMiddleware struct {
	session     *session.Session
	authService *service.AuthService
}

type AuthMiddlewareDependencies struct {
	fx.In

	Session     *session.Session
	AuthService *service.AuthService
}

func NewAuthMiddleware(deps AuthMiddlewareDependencies) *AuthMiddleware {
	return &AuthMiddleware{
		session:     deps.Session,
		authService: deps.AuthService,
	}
}

func (a *AuthMiddleware) Authenticated(ctx *fiber.Ctx) error {
	if err := a.authService.Authenticated(ctx); err != nil {
		if errors.Is(err, fiber.ErrUnauthorized) && strings.HasPrefix(ctx.Path(), "/_/") {
			return ctx.Redirect("/_/auth/login")
		}

		return err
	}

	return ctx.Next()
}
