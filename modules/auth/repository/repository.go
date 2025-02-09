package repository

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/core/cache"
	"github.com/mrrizkin/pohara/modules/database"
	"github.com/mrrizkin/pohara/modules/database/utils"
)

type AuthRepository struct {
	db    *database.Database
	cache *cache.Cache
}

type AuthRepositoryDependencies struct {
	fx.In

	DB    *database.Database
	Cache *cache.Cache
}

func NewAuthRepository(deps AuthRepositoryDependencies) *AuthRepository {
	return &AuthRepository{
		db:    deps.DB,
		cache: deps.Cache,
	}
}

type UserContext struct {
	User          *model.MUser
	UserAttribute *model.MUserAttribute
	UserSetting   *model.MUserSetting
	Roles         []model.MRole
}

func (a *AuthRepository) GetUserContext(ctx context.Context, userID uint) (*UserContext, error) {
	cacheKey := a.getUserContextCacheKey(userID)

	if val, found := a.cache.Get(ctx, cacheKey); found {
		return val.(*UserContext), nil
	}

	user, err := a.GetUser(userID)
	if err != nil {
		return nil, err
	}

	userAttribute, err := a.GetUserAttributes(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	userSetting, err := a.GetUserSettings(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	roles, err := a.GetUserRoles(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	roleIds := utils.Pluck(roles, func(mr model.MRole) uint { return mr.ID })
	rolePolicies, err := a.GetRolePolicies(roleIds...)
	if err != nil {
		return nil, err
	}

	rolePoliciesMap := utils.KeyByGroup(
		rolePolicies,
		func(rp RoleCfgPolicy) uint { return rp.RoleID },
	)

	for i, role := range roles {
		if rolePolicies, ok := rolePoliciesMap[role.ID]; ok {
			policies := utils.Map(rolePolicies, func(rp RoleCfgPolicy) model.CfgPolicy {
				return model.CfgPolicy{
					ID:        rp.ID,
					Name:      rp.Name,
					Condition: rp.Condition,
					Action:    rp.Action,
					Effect:    rp.Effect,
					Resource:  rp.Resource,
				}
			})
			roles[i].Policies = append(roles[i].Policies, policies...)
		}
	}

	userContext := &UserContext{
		User:          user,
		UserAttribute: userAttribute,
		UserSetting:   userSetting,
		Roles:         roles,
	}

	_ = a.cache.Set(ctx, cacheKey, userContext)
	return userContext, nil
}

func (a *AuthRepository) GetUser(userID uint) (*model.MUser, error) {
	var user model.MUser
	err := a.db.First(&user, userID).Error
	return &user, err
}

func (a *AuthRepository) GetUserAttributes(userID uint) (*model.MUserAttribute, error) {
	userAttributes := model.MUserAttribute{}
	if err := a.db.Where("user_id = ?", userID).First(&userAttributes).Error; err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &userAttributes, nil
}

func (a *AuthRepository) GetUserSettings(userID uint) (*model.MUserSetting, error) {
	userSetting := model.MUserSetting{}
	if err := a.db.Where("user_id = ?", userID).First(&userSetting).Error; err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &userSetting, nil
}

func (a *AuthRepository) GetUserRoles(userID uint) ([]model.MRole, error) {
	roles := make([]model.MRole, 0)
	if err := a.db.Table("m_role mr").
		Joins("INNER JOIN jt_user_role jur ON mr.id = jur.role_id").
		Where("jur.user_id = ?", userID).
		Select("mr.*").
		Find(&roles).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return roles, nil
}

type RoleCfgPolicy struct {
	model.CfgPolicy
	RoleID uint `json:"role_id"`
}

func (a *AuthRepository) GetRolePolicies(roleIds ...uint) ([]RoleCfgPolicy, error) {
	policies := make([]RoleCfgPolicy, 0)
	if err := a.db.Table("cfg_policy cp").
		Select("jrp.role_id, cp.*").
		Joins("INNER JOIN jt_role_policy jrp ON cp.id = jrp.policy_id").
		Where("jrp.role_id in (?)", roleIds).
		Find(&policies).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return policies, nil
}

func (*AuthRepository) DeleteUserContextCache(userID uint) error {
	return nil
}

func (*AuthRepository) getUserContextCacheKey(userID uint) string {
	return fmt.Sprintf("user:%d:context", userID)
}
