package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/validator"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"
)

type SetupController struct {
	log       *logger.ZeroLog
	inertia   *inertia.Inertia
	validator *validator.Validator
}

type SetupControllerDependencies struct {
	fx.In

	Logger    *logger.ZeroLog
	Inertia   *inertia.Inertia
	Validator *validator.Validator
}

func NewSetupController(deps SetupControllerDependencies) *SetupController {
	return &SetupController{
		log:       deps.Logger.Scope("setup_controller"),
		inertia:   deps.Inertia,
		validator: deps.Validator,
	}
}

func (c *SetupController) Index(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "setup/index", gonertia.Props{})
}
