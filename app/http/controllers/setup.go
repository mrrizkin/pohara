package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/common/sql"
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
	hash      *hash.Hashing

	userRepo *repository.UserRepository
}

type SetupControllerDependencies struct {
	fx.In

	Logger    *logger.ZeroLog
	Inertia   *inertia.Inertia
	Validator *validator.Validator
	Hashing   *hash.Hashing

	UserRepository *repository.UserRepository
}

func NewSetupController(deps SetupControllerDependencies) *SetupController {
	return &SetupController{
		log:       deps.Logger.Scope("setup_controller"),
		inertia:   deps.Inertia,
		validator: deps.Validator,
		hash:      deps.Hashing,

		userRepo: deps.UserRepository,
	}
}

func (c *SetupController) Index(ctx *fiber.Ctx) error {
	return c.inertia.Render(ctx, "setup/index", gonertia.Props{})
}

type AdminUserPayload struct {
	Username             string `json:"username"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type EmailSetupPayload struct {
	Driver     string `json:"driver"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"`
	Domain     string `json:"domain"`
	Secret     string `json:"secret"`
	Key        string `json:"key"`
	Token      string `json:"token"`
	APIKey     string `json:"api_key"`
}

type SiteSetupPayload struct {
	Domain      string `json:"domain"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Timezone    string `json:"timezone"`
	Currency    string `json:"currency"`
	Locale      string `json:"locale"`
}

type SetupPayload struct {
	AdminUser AdminUserPayload  `json:"admin_user"`
	Email     EmailSetupPayload `json:"email"`
	Site      SiteSetupPayload  `json:"site"`
}

func (c *SetupController) Setup(ctx *fiber.Ctx) error {
	var payload SetupPayload
	if err := c.validator.ParseBodyAndValidate(ctx, &payload); err != nil {
		return err
	}

	password, err := c.hash.Generate(payload.AdminUser.Password)
	if err != nil {
		return err
	}

	user := &model.MUser{
		Name:     payload.AdminUser.Name,
		Username: payload.AdminUser.Username,
		Email:    sql.String(payload.AdminUser.Email),
		Password: password,
	}

	if err := c.userRepo.Create(user); err != nil {
		return err
	}

	return c.inertia.Redirect(ctx, "/", fiber.StatusOK)
}
