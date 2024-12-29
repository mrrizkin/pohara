package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/common/validator"
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/mrrizkin/pohara/internal/web/template"
)

type Handler = func(*Ctx) error
type Map map[string]interface{}

type Router struct {
	core      fiber.Router
	template  *template.Template
	inertia   *inertia.Inertia
	validator *validator.Validator
	config    *config.App
	cache     ports.Cache
	log       ports.Logger
}

type Ctx struct {
	*fiber.Ctx
	router *Router
}
