package repository

import (
	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/database/db"
	"go.uber.org/fx"
)

type RoleRepository struct {
	db *db.Database
}

type RoleRepositoryDependencies struct {
	fx.In

	Database *db.Database
}

func NewRoleRepository(deps RoleRepositoryDependencies) *RoleRepository {
	return &RoleRepository{
		db: deps.Database,
	}
}

func (a *RoleRepository) CreateMRole(mRole *model.MRole) error {
	result := a.db.Create(mRole)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *RoleRepository) GetMRoleByID(id uint) (*model.MRole, error) {
	var mRole model.MRole
	result := a.db.First(&mRole, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &mRole, nil
}

func (r *RoleRepository) Find(
	search sql.StringNullable,
	paginateParams QueryPaginateParams,
) (*PaginationResult[model.MRole], error) {
	var roles []model.MRole
	query := r.db.Model(&roles)
	if search.Valid {
		query.Where("name ILIKE ?", "%"+search.String+"%")
	}
	return QueryPaginate(query, roles, paginateParams)
}

func (a *RoleRepository) GetAllMRoles() ([]model.MRole, error) {
	var mRoles []model.MRole
	result := a.db.Find(&mRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return mRoles, nil
}

func (a *RoleRepository) UpdateMRole(mRole *model.MRole) error {
	result := a.db.Save(mRole)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *RoleRepository) DeleteMRole(id uint) error {
	result := a.db.Delete(&model.MRole{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
