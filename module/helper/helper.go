package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/module/validator"
	"github.com/mrrizkin/pohara/module/view"
)

func ParseBodyAndValidate(ctx *fiber.Ctx, validator *validator.Validator, out interface{}) error {
	err := ctx.BodyParser(out)
	if err != nil {
		return &fiber.Error{
			Code:    400,
			Message: "payload not valid",
		}
	}

	err = validator.MustValidate(out)
	if err != nil {
		return err
	}

	return nil
}

func Render(ctx *fiber.Ctx, view *view.View, template string, data fiber.Map) error {
	html, err := view.Render(template, data)
	if err != nil {
		return err
	}

	return ctx.Type("html").Send(html)
}
