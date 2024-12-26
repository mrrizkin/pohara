package user

import (
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/models"
	"github.com/mrrizkin/pohara/module/logger"
	"github.com/mrrizkin/pohara/module/server"
	"github.com/mrrizkin/pohara/module/validator"
	"go.uber.org/fx"
)

type UserHandler struct {
	log         *logger.Logger
	validator   *validator.Validator
	userService *UserService
}

type HandlerDependencies struct {
	fx.In

	Log         *logger.Logger
	Validator   *validator.Validator
	UserService *UserService
}

type HandlerResult struct {
	fx.Out

	UserHandler *UserHandler
}

func Handler(deps HandlerDependencies) HandlerResult {
	return HandlerResult{
		UserHandler: &UserHandler{
			log:         deps.Log.Scope("user_handler"),
			validator:   deps.Validator,
			userService: deps.UserService,
		},
	}
}

// UserCreate godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided information
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User							true	"User information"
//	@Success		200		{object}	fiber.Map{data=models.User}	"Successfully created user"
//	@Failure		400		{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		500		{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user [post]
func (h *UserHandler) UserCreate(ctx *server.Ctx) error {
	payload := new(models.User)
	err := ctx.ParseBodyAndValidate(payload)
	if err != nil {
		h.log.Error("failed to parse and validate payload", "error", err)
		return err
	}

	user, err := h.userService.Create(payload)
	if err != nil {
		h.log.Error("failed to create user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to create user: %s", err),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "user created successfully",
		"data":    user,
	})
}

// UserFindAll godoc
//
//	@Summary		Get all users
//	@Description	Retrieve a list of all users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"
//	@Param			per_page	query		int																false	"Number of items per page"
//	@Success		200			{object}	fiber.Map{data=[]models.User,meta=fiber.Map}	"Successfully retrieved users"
//	@Failure		500			{object}	validator.GlobalErrorResponse									"Internal server error"
//	@Router			/user [get]
func (h *UserHandler) UserFindAll(ctx *server.Ctx) error {
	page := ctx.QueryInt("page", 1)
	perPage := ctx.QueryInt("per_page", 10)

	users, err := h.userService.FindAll(page, perPage)
	if err != nil {
		h.log.Error("failed to get users", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get users: %s", err),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "users retrieved successfully",
		"data":    users["result"],
		"meta": fiber.Map{
			"page":       page,
			"per_page":   perPage,
			"total":      users["total"],
			"page_count": math.Ceil(float64(users["total"].(int)) / float64(perPage)),
		},
	})
}

// UserFindByID godoc
//
//	@Summary		Get a user by ID
//	@Description	Retrieve a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int									true	"User ID"
//	@Success		200	{object}	fiber.Map{data=models.User}	"Successfully retrieved user"
//	@Failure		400	{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		404	{object}	validator.GlobalErrorResponse		"User not found"
//	@Failure		500	{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user/{id} [get]
func (h *UserHandler) UserFindByID(ctx *server.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		h.log.Error("failed to parse id", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	user, err := h.userService.FindByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "user not found",
			}
		}

		h.log.Error("failed to get user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get user: %s", err),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "user retrieved successfully",
		"data":    user,
	})
}

// UserUpdate godoc
//
//	@Summary		Update a user
//	@Description	Update a user's information by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"User ID"
//	@Param			user	body		models.User							true	"Updated user information"
//	@Success		200		{object}	fiber.Map{data=models.User}	"Successfully updated user"
//	@Failure		400		{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		500		{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user/{id} [put]
func (h *UserHandler) UserUpdate(ctx *server.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		h.log.Error("failed to parse id", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	payload := new(models.User)
	err = ctx.ParseBodyAndValidate(payload)
	if err != nil {
		h.log.Error("failed to parse and validate payload", "error", err)
		return err
	}

	user, err := h.userService.Update(uint(id), payload)
	if err != nil {
		h.log.Error("failed to update user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to update user: %s", err),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "user updated successfully",
		"data":    user,
	})
}

// UserDelete godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"User ID"
//	@Success		200	{object}	fiber.Map					"Successfully deleted user"
//	@Failure		400	{object}	validator.GlobalErrorResponse	"Bad request"
//	@Failure		401	{object}	validator.GlobalErrorResponse	"Unauthorized"
//	@Failure		500	{object}	validator.GlobalErrorResponse	"Internal server error"
//	@Router			/user/{id} [delete]
func (c *UserHandler) UserDelete(ctx *server.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Error("failed to parse id", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	err = c.userService.Delete(uint(id))
	if err != nil {
		c.log.Error("failed to delete user", "err", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to delete user: %s", err),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted successfully",
	})
}
