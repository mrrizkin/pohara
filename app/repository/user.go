package repository

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/action"
	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/abac/access"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/database/db"
	"github.com/mrrizkin/pohara/modules/logger"
	"github.com/mrrizkin/pohara/modules/validator"
)

type UserRepository struct {
	db        *db.Database
	log       *logger.Logger
	validator *validator.Validator
	hashing   *hash.Hashing
}

type UserRepositoryDependencies struct {
	fx.In

	Database  *db.Database
	Logger    *logger.Logger
	Validator *validator.Validator
	Hashing   *hash.Hashing
}

func NewUserRepository(deps UserRepositoryDependencies) *UserRepository {
	return &UserRepository{
		log:       deps.Logger.System().Scope("user_repository"),
		validator: deps.Validator,
		hashing:   deps.Hashing,
		db:        deps.Database,
	}
}

func (r *UserRepository) SetupSuperUser(user *model.MUser) error {
	policy := model.CfgPolicy{
		Name:     "Allow All Function",
		Action:   action.SpecialAll,
		Effect:   access.EffectAllow,
		Resource: sql.String("all"),
	}

	role := model.MRole{
		Name:        "Super User",
		Description: "the most highly privilege role",
	}

	tx := r.db.Begin()

	if err := tx.FirstOrCreate(&policy).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.FirstOrCreate(&role).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.FirstOrCreate(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.FirstOrCreate(&model.JtUserRole{
		RoleID: role.ID,
		UserID: user.ID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.FirstOrCreate(&model.JtRolePolicy{
		PolicyID: policy.ID,
		RoleID:   role.ID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *UserRepository) Create(user *model.MUser) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Find(
	search sql.StringNullable,
	paginateParams QueryPaginateParams,
) (*PaginationResult[model.MUser], error) {
	var users []model.MUser
	query := r.db.Model(&users)
	if search.Valid {
		query.Where("name ILIKE ?", "%"+search.String+"%")
	}
	return QueryPaginate(query, users, paginateParams)
}

func (r *UserRepository) FindByID(id uint) (*model.MUser, error) {
	var user model.MUser
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*model.MUser, error) {
	var user model.MUser
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.MUser) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.MUser{}, id).Error
}
