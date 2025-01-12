package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/validator"
	"go.uber.org/fx"
)

type UserController struct {
	log       *logger.ZeroLog
	validator *validator.Validator
	userRepo  *repository.UserRepository
	hashing   *hash.Hashing
}

type UserControllerDependencies struct {
	fx.In

	Logger    *logger.ZeroLog
	Validator *validator.Validator
	Hashing   *hash.Hashing

	UserRepository *repository.UserRepository
}

func NewUserController(deps UserControllerDependencies) *UserController {
	return &UserController{
		log:       deps.Logger.Scope("user_controller"),
		validator: deps.Validator,
		hashing:   deps.Hashing,
		userRepo:  deps.UserRepository,
	}
}

type UserCreatePayload struct {
	Name     string             `json:"name"     validate:"required"`
	Username string             `json:"username" validate:"required"`
	Password string             `json:"password" validate:"required"`
	Email    sql.StringNullable `json:"email"`
}

type UserUpdatePayload struct {
	Name     string             `json:"name"`
	Username sql.StringNullable `json:"username"`
	Password sql.StringNullable `json:"password"`
	Email    sql.StringNullable `json:"email"`
}

// UserCreate godoc
//
//	@Summary      Create a new user
//	@Description  Create a new user with the provided information
//	@Tags         Users
//	@Accept       json
//	@Produce		  json
//	@Param			  user	    body		  UserCreatePayload	true	"New user information"
//	@Success		  200		    {object}	fiber.Map{status=string,message=string}	"Successfully created user"
//	@Failure		  400		    {object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		  500		    {object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			  /user     [post]
func (c *UserController) UserCreate(ctx *fiber.Ctx) error {
	var err error

	var payload UserCreatePayload
	err = c.validator.ParseBodyAndValidate(ctx, &payload)
	if err != nil {
		cause := "error parse and validate"
		c.log.Error(cause, "error", err)
		return err
	}

	if payload.Password != "" {
		cause := "password is required"
		c.log.Error(cause)
		return &fiber.Error{
			Code:    fiber.StatusBadGateway,
			Message: fmt.Sprintf("failed to create user: %s", cause),
		}
	}

	hash, err := c.hashing.Generate(payload.Password)
	if err != nil {
		cause := "error hashing password"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to create user: %s", cause),
		}
	}

	err = c.userRepo.Create(&model.MUser{
		Name:     payload.Name,
		Username: payload.Username,
		Password: hash,
		Email:    payload.Email,
	})
	if err != nil {
		cause := "error create user to database"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to create user: %s", cause),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "user created successfully",
	})
}

// UserFind godoc
//
//	@Summary		Get all users
//	@Description	Retrieve a list of all users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"
//	@Param			per_page	query		int																false	"Number of items per page"
//	@Success		200			{object}	fiber.Map{status=string,message=string,data=[]model.User,meta=fiber.Map{page=int,limit=int,total=int,total_page=int}}	"Successfully retrieved users"
//	@Failure		500			{object}	validator.GlobalErrorResponse									"Internal server error"
//	@Router			/user [get]
func (c *UserController) UserFind(ctx *fiber.Ctx) error {
	page := int64(ctx.QueryInt("page", 1))
	limit := int64(ctx.QueryInt("limit", 10))
	searchQ := ctx.Query("q", "")

	search := sql.StringNull()
	if searchQ != "" {
		search.Valid = true
		search.String = searchQ
	}

	result, err := c.userRepo.Find(search, sql.Int64(page), sql.Int64(limit))
	if err != nil {
		cause := "error find users"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get users: %s", cause),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "users retrieved successfully",
		"data":    result.Data,
		"meta": fiber.Map{
			"page":       result.Page,
			"limit":      result.Limit,
			"total":      result.Total,
			"total_page": result.TotalPage,
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
//	@Success		200	{object}	fiber.Map{status=string,message=string,data=model.User}	"Successfully retrieved user"
//	@Failure		400	{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		404	{object}	validator.GlobalErrorResponse		"User not found"
//	@Failure		500	{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user/{id} [get]
func (c *UserController) UserFindByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		cause := "error parse id required"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			cause := "user not found"
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: cause,
			}
		}

		cause := "error get user"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get user: %s", cause),
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
//	@Param			user	body		UserUpdatePayload	true	"Updated user information"
//	@Success		200		{object}	fiber.Map{status=string,message=string,data=model.User}	"Successfully updated user"
//	@Failure		400		{object}	validator.GlobalErrorResponse		"Bad request"
//	@Failure		500		{object}	validator.GlobalErrorResponse		"Internal server error"
//	@Router			/user/{id} [put]
func (c *UserController) UserUpdate(ctx *fiber.Ctx) error {
	var err error
	var payload UserUpdatePayload

	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Error("failed to parse id", "error", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	err = c.validator.ParseBodyAndValidate(ctx, &payload)
	if err != nil {
		c.log.Error("failed to parse and validate payload", "error", err)
		return err
	}

	user, err := c.userRepo.FindByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			cause := "user not found"
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: cause,
			}
		}

		cause := "error get user"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to get user: %s", cause),
		}
	}

	if payload.Password.Valid {
		hash, err := c.hashing.Generate(payload.Password.String)
		if err != nil {
			cause := "error hashing password"
			c.log.Error(cause, "error", err)
			return &fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: fmt.Sprintf("failed to create user: %s", cause),
			}
		}

		user.Password = hash
	}

	user.Name = payload.Name
	user.Email = payload.Email
	if payload.Username.Valid {
		user.Username = payload.Username.String
	}

	err = c.userRepo.Update(user)
	if err != nil {
		cause := "error update user to database"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to update user: %s", cause),
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
//	@Success		200	{object}	fiber.Map{status=string,message=string}	"Successfully deleted user"
//	@Failure		400	{object}	validator.GlobalErrorResponse	"Bad request"
//	@Failure		401	{object}	validator.GlobalErrorResponse	"Unauthorized"
//	@Failure		500	{object}	validator.GlobalErrorResponse	"Internal server error"
//	@Router			/user/{id} [delete]
func (c *UserController) UserDelete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Error("failed to parse id", "error", err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid id",
		}
	}

	err = c.userRepo.Delete(uint(id))
	if err != nil {
		cause := "error delete user in database"
		c.log.Error(cause, "error", err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: fmt.Sprintf("failed to delete user: %s", cause),
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted successfully",
	})
}
