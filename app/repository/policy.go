package repository

import (
	"errors"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/database/db"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type PolicyRepository struct {
	db *db.Database
}

type PolicyRepositoryDependencies struct {
	fx.In

	Database *db.Database
}

func NewPolicyRepository(deps PolicyRepositoryDependencies) *PolicyRepository {
	return &PolicyRepository{
		db: deps.Database,
	}
}

func (a *PolicyRepository) CreateMPolicy(cfgPolicy *model.CfgPolicy) error {
	result := a.db.Create(cfgPolicy)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *PolicyRepository) GetMPolicyByID(id uint) (*model.CfgPolicy, error) {
	var cfgPolicy model.CfgPolicy
	result := a.db.First(&cfgPolicy, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cfgPolicy, nil
}

func (a *PolicyRepository) GetPolicyRoles(id uint) ([]model.MRole, error) {
	roles := make([]model.MRole, 0)
	err := a.db.Table("m_role mr").
		Select("mr.*").
		Joins("INNER JOIN jt_role_policy jrp ON mr.id = jrp.role_id").
		Where("jrp.policy_id = ?", id).
		Find(&roles).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return roles, err
}

func (r *PolicyRepository) Find(
	search sql.StringNullable,
	paginateParams QueryPaginateParams,
) (*PaginationResult[model.CfgPolicy], error) {
	var policies []model.CfgPolicy
	query := r.db.Model(&policies)
	if search.Valid {
		query.Where("name LIKE ?", "%"+search.String+"%")
	}
	return QueryPaginate(query, policies, paginateParams)
}

func (a *PolicyRepository) GetAllMPolicys() ([]model.CfgPolicy, error) {
	var policies []model.CfgPolicy
	result := a.db.Find(&policies)
	if result.Error != nil {
		return nil, result.Error
	}
	return policies, nil
}

func (a *PolicyRepository) UpdateMPolicy(cfgPolicy *model.CfgPolicy) error {
	result := a.db.Save(cfgPolicy)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *PolicyRepository) DeleteMPolicy(id uint) error {
	result := a.db.Delete(&model.CfgPolicy{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
