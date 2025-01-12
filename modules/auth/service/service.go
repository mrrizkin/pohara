package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/auth/access"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Authenticated(c *fiber.Ctx) error {
	return c.Next()
}

func (a *AuthService) Login() error {
	return nil
}

func (a *AuthService) Logout() error {
	return nil
}

func (a *AuthService) Can(
	ctx *fiber.Ctx,
	action access.Action,
	resources interface{},
) bool {
	return false
}
