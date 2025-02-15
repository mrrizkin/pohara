package repository

import (
	"errors"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/database/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RolePolicy struct {
	model.CfgPolicy
	RoleID uint `json:"role_id"`
}

type ListRolePolicy []RolePolicy

func (a *AuthRepository) CreateMRole(mRole *model.MRole) error {
	result := a.db.Create(mRole)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *AuthRepository) GetMRoleByID(id uint) (*model.MRole, error) {
	var mRole model.MRole
	result := a.db.First(&mRole, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &mRole, nil
}

func (a *AuthRepository) GetAllMRoles() ([]model.MRole, error) {
	var mRoles []model.MRole
	result := a.db.Find(&mRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return mRoles, nil
}

func (a *AuthRepository) UpdateMRole(mRole *model.MRole) error {
	result := a.db.Save(mRole)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *AuthRepository) DeleteMRole(id uint) error {
	result := a.db.Delete(&model.MRole{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *AuthRepository) ReconcilRoleToPolicy(
	role *model.MRole,
	policies ...model.CfgPolicy,
) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Batch insert new role-policy assignments (skip if already exists)
		var jtRolePolicies []model.JtRolePolicy
		for _, policy := range policies {
			jtRolePolicies = append(jtRolePolicies, model.JtRolePolicy{
				RoleID:   role.ID,
				PolicyID: policy.ID,
			})
		}

		// Batch insert with conflict handling
		if len(jtRolePolicies) > 0 {
			result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&jtRolePolicies)
			if result.Error != nil {
				return result.Error
			}
		}

		// Step 2: Delete role-policy assignments that are no longer in the updated list
		listPolicyId := utils.Pluck(policies, func(p model.CfgPolicy) uint { return p.ID })
		if len(listPolicyId) > 0 {
			result := tx.Where("role_id = ? AND policy_id NOT IN (?)", role.ID, listPolicyId).
				Delete(&model.JtRolePolicy{})
			if result.Error != nil {
				return result.Error
			}
		} else {
			// If no policies are provided, delete all assignments for the role
			result := tx.Where("role_id = ?", role.ID).
				Delete(&model.JtRolePolicy{})
			if result.Error != nil {
				return result.Error
			}
		}

		return nil
	})
}

func (a *AuthRepository) getRolePolicies(roleIds ...uint) (ListRolePolicy, error) {
	policies := make(ListRolePolicy, 0)
	err := a.db.Table("cfg_policy cp").
		Select("jrp.role_id, cp.*").
		Joins("INNER JOIN jt_role_policy jrp ON cp.id = jrp.policy_id").
		Where("jrp.role_id IN (?)", roleIds).
		Find(&policies).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return policies, nil
}

func (lrp ListRolePolicy) mapToRole(roles ...model.MRole) {
	lrpMap := utils.KeyByGroup(lrp, func(rp RolePolicy) uint {
		return rp.RoleID
	})

	for i, role := range roles {
		if rolePolicies, ok := lrpMap[role.ID]; ok {
			policies := utils.Map(rolePolicies, func(rp RolePolicy) model.CfgPolicy {
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
}
