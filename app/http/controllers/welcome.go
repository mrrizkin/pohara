package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/neoweb/templ"
	"github.com/mrrizkin/pohara/resources/views/pages"
)

type WelcomeController struct {
}

func NewWelcomeController() *WelcomeController {
	return &WelcomeController{}
}

func (c *WelcomeController) Index(ctx *fiber.Ctx) error {
	return templ.Render(ctx, pages.Welcome())
}
