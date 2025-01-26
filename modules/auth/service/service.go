package service

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/modules/auth/access"
	"github.com/mrrizkin/pohara/modules/auth/repository"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/session"
)

type AuthService struct {
	sessionStore *session.Store
	log          *logger.ZeroLog
	authRepo     *repository.AuthRepository
}

type AuthServiceDependencies struct {
	fx.In

	SessionStore *session.Store
	Logger       *logger.ZeroLog
	AuthRepo     *repository.AuthRepository
}

func NewAuthService(deps AuthServiceDependencies) *AuthService {
	return &AuthService{
		sessionStore: deps.SessionStore,
		log:          deps.Logger,
		authRepo:     deps.AuthRepo,
	}
}

func (a *AuthService) Authenticated(ctx *fiber.Ctx) error {
	sess, err := a.sessionStore.Get(ctx)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	authUID := sess.Get("__auth-uid")
	if authUID == nil {
		return fiber.ErrUnauthorized
	}

	uid, ok := authUID.(uint)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, err := a.authRepo.GetUser(uid)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	userAttributes, err := a.authRepo.GetUserAttributes(uid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.log.Error("failed to get user attributes", "error", err)
		return fiber.ErrUnauthorized
	}

	ctx.Locals("__auth-user", user)
	ctx.Locals("__auth-user-attributes", userAttributes)

	return nil
}

func (a *AuthService) Login(ctx *fiber.Ctx, uid uint) error {
	sess, err := a.sessionStore.Get(ctx)
	if err != nil {
		return err
	}

	sess.Set("__auth-uid", uid)
	err = sess.Save()
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Logout(ctx *fiber.Ctx) error {
	sess, err := a.sessionStore.Get(ctx)
	if err != nil {
		return err
	}

	sess.Delete("__auth-uid")
	err = sess.Save()
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Can(
	ctx *fiber.Ctx,
	action access.Action,
	resources interface{},
) bool {
	return false
}
