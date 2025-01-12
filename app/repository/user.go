package repository

import (
	"math"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/common/hash"
	"github.com/mrrizkin/pohara/modules/common/sql"
	"github.com/mrrizkin/pohara/modules/core/database"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/core/validator"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UserRepository struct {
	db        *database.GormDB
	log       *logger.ZeroLog
	validator *validator.Validator
	hashing   *hash.Hashing
}

type UserRepositoryDependencies struct {
	fx.In

	Database  *database.GormDB
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
	return r.db.Create(user).Error
}

func (r *UserRepository) Find(
	search sql.StringNullable,
	page, limit sql.Int64Nullable,
) (result *sql.PaginationResult, err error) {
	var users []model.MUser
	var total int64

	query := r.db.Model(&users)
	err = query.Session(&gorm.Session{NewDB: true}).Count(&total).Error
	if err != nil {
		return
	}

	if limit.Valid {
		query = query.Limit(int(limit.Int64))
	}

	offset := sql.Int64Null()
	if page.Valid && page.Int64 != 0 {
		offset.Valid = true
		offset.Int64 = (page.Int64 - 1)
		if limit.Int64 != 0 {
			offset.Int64 = page.Int64 * limit.Int64
		}

		query = query.Offset(int(offset.Int64))
	}

	err = query.Find(&users).Error
	if err != nil {
		return
	}

	totalPage := sql.Int64Null()
	if offset.Valid && limit.Valid && limit.Int64 != 0 {
		totalPage.Valid = true
		totalPage.Int64 = int64(math.Ceil(float64(offset.Int64)/float64(limit.Int64))) + 1
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
