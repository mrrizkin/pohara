package repository

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/model"
	"github.com/mrrizkin/pohara/modules/database"
)

type AuthRepository struct {
	db *database.Database
}

type AuthRepositoryDependencies struct {
	fx.In
	DB *database.Database
}

func NewAuthRepository(deps AuthRepositoryDependencies) *AuthRepository {
	return &AuthRepository{
		db: deps.DB,
	}
}

func (a *AuthRepository) GetUser(uid uint) (*model.MUser, error) {
	var user model.MUser
	err := a.db.Builder().Table("m_user").Select().Where("id = ?", uid).Get(&user)
	return &user, err
}

func (a *AuthRepository) GetUserAttributes(uid uint) (*model.MUserAttribute, error) {
	var userAttributes model.MUserAttribute
	err := a.db.Builder().Table("m_user_attribute").Select().Where("user_id = ?", uid).Get(&userAttributes)
	return &userAttributes, err
}
