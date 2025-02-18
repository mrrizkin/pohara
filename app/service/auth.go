package service

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/modules/abac"
	"github.com/mrrizkin/pohara/modules/abac/access"
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/mrrizkin/pohara/modules/session"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type AuthService struct {
	sessionStore *session.Store
	log          *logger.Logger
	authRepo     *repository.AuthRepository
}

type AuthServiceDependencies struct {
	fx.In

	SessionStore *session.Store
	Logger       *logger.Logger
	AuthRepo     *repository.AuthRepository
}

func NewAuthService(deps AuthServiceDependencies) *AuthService {
	return &AuthService{
		sessionStore: deps.SessionStore,
		log:          deps.Logger,
		authRepo:     deps.AuthRepo,
	}
}

type Subject struct {
	*model.MUser
	Attribute *model.MUserAttribute
	Setting   *model.MUserSetting
	roles     []model.MRole
}

func (a *Subject) GetPolicies() ([]abac.Policy, error) {
	policies := make([]abac.Policy, 0)

	for _, role := range a.roles {
		for _, policy := range role.Policies {
			policies = append(policies, policy)
		}
	}

	return policies, nil
}

func (a AuthService) GetActionSpecialAll() access.Action {
	return action.SpecialAll
}

func (a AuthService) GetSubjectName() string {
	return "User"
}

func (a AuthService) GetSubject(ctx *fiber.Ctx) (abac.Subject, error) {
	userContext, ok := ctx.Locals("__auth-user-context").(*repository.UserContext)
	if !ok {
		return nil, errors.New("subject not found")
	}

	if userContext == nil {
		return nil, errors.New("subject not found")
	}

	return &Subject{
		MUser:     userContext.User,
		Attribute: userContext.UserAttribute,
		Setting:   userContext.UserSetting,
		roles:     userContext.Roles,
	}, nil
}

func (a AuthService) GetEnv(ctx *fiber.Ctx) (any, error) {
	return nil, nil
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

	userContext, err := a.authRepo.GetUserContext(ctx.Context(), uid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.log.Error("failed to get user attributes", "error", err)
		return fiber.ErrUnauthorized
	}

	ctx.Locals("__auth-user-context", userContext)
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
