package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/module/validator"
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
