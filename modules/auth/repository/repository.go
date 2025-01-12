package repository

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/core/database"
)

type AuthRepository struct {
	db *database.GormDB
}

type AuthRepositoryDependencies struct {
	fx.In
	DB *database.GormDB
}

func NewAuthRepository(deps AuthRepositoryDependencies) *AuthRepository {
	return &AuthRepository{
		db: deps.DB,
	}
}

func (a *AuthRepository) GetUser(uid uint) (*model.MUser, error) {
	var user model.MUser
	err := a.db.First(&user, uid).Error
	return &user, err
}

func (a *AuthRepository) GetUserAttributes(uid uint) (*model.MUserAttribute, error) {
	var userAttributes model.MUserAttribute
	err := a.db.Where("user_id = ?", uid).First(&userAttributes).Error
	return &userAttributes, err
}
