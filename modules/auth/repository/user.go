package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/database/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserContext struct {
	User          *model.MUser
	UserAttribute *model.MUserAttribute
	UserSetting   *model.MUserSetting
	Roles         []model.MRole
}

func (a *AuthRepository) UserModified(ctx context.Context, userID uint) error {
	return a.deleteUserContextCache(ctx, userID)
}

func (a *AuthRepository) GetUserContext(ctx context.Context, userID uint) (*UserContext, error) {
	if val, found := a.getUserContextCache(ctx, userID); found {
		return val, nil
	}

	user, err := a.getUser(userID)
	if err != nil {
		return nil, err
	}

	userAttribute, err := a.getUserAttributes(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	userSetting, err := a.getUserSettings(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	roles, err := a.getUserRoles(userID, true)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	userContext := UserContext{
		User:          user,
		UserAttribute: userAttribute,
		UserSetting:   userSetting,
		Roles:         roles,
	}

	err = a.createUserContextCache(ctx, userID, &userContext)
	return &userContext, err
}

func (a *AuthRepository) AssignUserToRole(user *model.MUser, roles ...model.MRole) error {
	// Prepare a slice of JtRolePolicy for batch insertion
	var jtUserRoles []model.JtUserRole
	for _, role := range roles {
		jtUserRoles = append(jtUserRoles, model.JtUserRole{
			UserID: user.ID,
			RoleID: role.ID,
		})
	}

	// Perform batch insertion
	result := a.db.CreateInBatches(jtUserRoles, len(jtUserRoles)) // Batch size = total number of policies
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a *AuthRepository) ReconcilUserToRole(user *model.MUser, roles ...model.MRole) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Batch insert new user-role assignments (skip if already exists)
		var jtUserRole []model.JtUserRole
		for _, role := range roles {
			jtUserRole = append(jtUserRole, model.JtUserRole{
				UserID: user.ID,
				RoleID: role.ID,
			})
		}

		// Batch insert with conflict handling
		if len(jtUserRole) > 0 {
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&jtUserRole)
			if result.Error != nil {
				return result.Error
			}
		}

		// Step 2: Delete user-role assignments that are no longer in the updated list
		listRoleId := utils.Pluck(roles, func(p model.MRole) uint { return p.ID })
		if len(listRoleId) > 0 {
			result := tx.Where("user_id = ? AND role_id NOT IN (?)", user.ID, listRoleId).
				Delete(&model.JtUserRole{})
			if result.Error != nil {
				return result.Error
			}
		} else {
			// If no rorles are provided, delete all assignments for the user
			result := tx.Where("user_id = ?", user.ID).
				Delete(&model.JtUserRole{})
			if result.Error != nil {
				return result.Error
			}
		}

		return nil
	})
}

func (a *AuthRepository) getUser(userID uint) (*model.MUser, error) {
	var user model.MUser
	err := a.db.First(&user, userID).Error
	return &user, err
}

func (a *AuthRepository) getUserAttributes(userID uint) (*model.MUserAttribute, error) {
	var userAttributes model.MUserAttribute
	err := a.db.Where("user_id = ?", userID).First(&userAttributes).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return &userAttributes, err
}

func (a *AuthRepository) getUserSettings(userID uint) (*model.MUserSetting, error) {
	var userSetting model.MUserSetting
	err := a.db.Where("user_id = ?", userID).First(&userSetting).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return &userSetting, err
}

func (a *AuthRepository) getUserRoles(userID uint, preload bool) ([]model.MRole, error) {
	roles := make([]model.MRole, 0)
	err := a.db.Table("m_role mr").
		Joins("INNER JOIN jt_user_role jur ON mr.id = jur.role_id").
		Where("jur.user_id = ?", userID).
		Select("mr.*").
		Find(&roles).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if preload {
		roleIds := utils.Pluck(roles, func(mr model.MRole) uint { return mr.ID })
		listRolePolicy, err := a.getRolePolicies(roleIds...)
		if err != nil {
			return nil, err
		}
		listRolePolicy.mapToRole(roles...)
	}

	return roles, err
}

func (a *AuthRepository) getUserContextCache(ctx context.Context, userID uint) (*UserContext, bool) {
	val, found := a.cache.Get(ctx, a.getUserContextCacheKey(userID))
	if !found {
		return nil, false
	}
	userContext, ok := val.(*UserContext)
	if !ok {
		return nil, false
	}
	return userContext, true
}

func (a *AuthRepository) createUserContextCache(ctx context.Context, userID uint, value *UserContext) error {
	return a.cache.Set(ctx, a.getUserContextCacheKey(userID), value)
}

func (a *AuthRepository) deleteUserContextCache(ctx context.Context, userID uint) error {
	return a.cache.Delete(ctx, a.getUserContextCacheKey(userID))
}

func (*AuthRepository) getUserContextCacheKey(userID uint) string {
	return fmt.Sprintf("user:%d:context", userID)
}
