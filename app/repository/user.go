package repository

import (
	"fmt"
	"math"

	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/validator"
	"github.com/mrrizkin/pohara/modules/database"
)

type UserRepository struct {
	db        *database.Database
	log       *logger.ZeroLog
	validator *validator.Validator
	hashing   *hash.Hashing
}

type UserRepositoryDependencies struct {
	fx.In

	Database  *database.Database
	Logger    *logger.ZeroLog
	Validator *validator.Validator
	Hashing   *hash.Hashing
}

func NewUserRepository(deps UserRepositoryDependencies) *UserRepository {
	return &UserRepository{
		log:       deps.Logger.Scope("user_repository"),
		validator: deps.Validator,
		hashing:   deps.Hashing,
		db:        deps.Database,
	}
}

func (r *UserRepository) Create(user *model.MUser) error {
	sqlx := r.db.Builder()
	_, err := sqlx.Table("m_user").Insert().Values(user).Exec()
	return err
}

func (r *UserRepository) Find(
	search sql.StringNullable,
	page, limit sql.Int64Nullable,
) (result *sql.PaginationResult, err error) {
	var users []model.MUser

	sqlx := r.db.Builder()
	sqlx.Table("m_user")
	total, err := sqlx.Clone().Count()
	if err != nil {
		return
	}

	if limit.Valid {
		sqlx.Limit(int(limit.Int64))
	}

	offset := sql.Int64Null()
	if page.Valid && page.Int64 != 0 {
		offset.Valid = true
		offset.Int64 = (page.Int64 - 1)
		if limit.Int64 != 0 {
			offset.Int64 = page.Int64 * limit.Int64
		}

		sqlx.Offset(int(offset.Int64))
	}

	err = sqlx.All(&users)
	if err != nil {
		return
	}

	totalPage := sql.Int64Null()
	if offset.Valid && limit.Valid && limit.Int64 != 0 {
		totalPage.Valid = true
		totalPage.Int64 = int64(math.Ceil(float64(total)/float64(limit.Int64))) + 1
	}

	result = &sql.PaginationResult{
		Data:      users,
		Total:     total,
		TotalPage: totalPage,
		Page:      page,
		Limit:     limit,
	}

	return
}

func (r *UserRepository) FindByID(id uint) (*model.MUser, error) {
	var user model.MUser

	sqlx := r.db.Builder()
	err := sqlx.Table("m_user").Select().Where("id = ?", id).Get(&user)
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*model.MUser, error) {
	var user model.MUser

	sqlx := r.db.Builder()
	err := sqlx.Table("m_user").Select().Where("email = ?", email).Get(&user)
	return &user, err
}

func (r *UserRepository) Update(user *model.MUser) error {
	sqlx := r.db.Builder()
	_, err := sqlx.Table("m_user").Update().
		Set("name", user.Name).
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Where("id = ?", user.ID).
		Exec()

	return err
}

func (r *UserRepository) Delete(id uint) error {
	sqlx := r.db.Builder()
	result, err := sqlx.Table("m_user").Delete().Where("id = ?", id).Exec()
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return fmt.Errorf("nothing deleted")
	}

	return nil
}
