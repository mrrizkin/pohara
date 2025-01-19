package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/romsar/gonertia"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/modules/auth/service"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/validator"
	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
)

type AuthController struct {
	log       *logger.ZeroLog
	inertia   *inertia.Inertia
	validator *validator.Validator
	hash      *hash.Hashing

	authService *service.AuthService
	userRepo    *repository.UserRepository
}

type AuthControllerDependencies struct {
	fx.In

	Logger    *logger.ZeroLog
	Inertia   *inertia.Inertia
	Validator *validator.Validator
	Hashing   *hash.Hashing

	AuthService    *service.AuthService
	UserRepository *repository.UserRepository
}

func NewAuthController(deps AuthControllerDependencies) *AuthController {
	return &AuthController{
		log:       deps.Logger.Scope("auth_controller"),
		inertia:   deps.Inertia,
		validator: deps.Validator,
		hash:      deps.Hashing,

		authService: deps.AuthService,
		userRepo:    deps.UserRepository,
	}
}

func (c *AuthController) LoginPage(ctx *fiber.Ctx) error {
	if err := c.authService.Authenticated(ctx); err == nil {
		return ctx.Redirect("/_/")
	}

	return c.inertia.Render(ctx, "auth/login", gonertia.Props{})
}

func (c *AuthController) RegisterPage(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "auth/register", gonertia.Props{})
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var input LoginInput
	if err := c.validator.ParseBodyAndValidate(ctx, &input); err != nil {
		cause := "error parse and validate"
		c.log.Error(cause, "error", err)
		return err
	}

	// get the user by the username
	user, err := c.userRepo.FindByEmail(input.Email)
	if err != nil {
		cause := "email or password is incorrect"
		c.log.Error(cause, "error", err)
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("failed: %s", cause))
	}

	// compore the password
	if match, err := c.hash.Compare(input.Password, user.Password); !match || err != nil {
		cause := "email or password is incorrect"
		c.log.Error(cause, "error", err)
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("failed: %s", cause))
	}

	// login
	if err := c.authService.Login(ctx, user.ID); err != nil {
		cause := "error login"
		c.log.Error(cause, "error", err)
		return err
	}

	return c.inertia.Redirect(ctx, "/_/")
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	return nil
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	if err := c.authService.Logout(ctx); err != nil {
		cause := "error logout"
		c.log.Error(cause, "error", err)
		return err
	}

	c.inertia.ClearHistory(ctx)
	return c.inertia.Redirect(ctx, "/_/auth/login")
}
