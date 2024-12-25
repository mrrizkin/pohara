package user

import (
	"github.com/mrrizkin/pohara/models"
	"github.com/mrrizkin/pohara/module/database"
	"go.uber.org/fx"
)

type UserRepository struct {
	db *database.Database
}

type RepositoryDependencies struct {
	fx.In

	Db *database.Database
}

type RepositoryResult struct {
	fx.Out

	UserRepository *UserRepository
}

func Repository(deps RepositoryDependencies) RepositoryResult {
	return RepositoryResult{
		UserRepository: &UserRepository{
			db: deps.Db,
		},
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindAll(
	page int,
	perPage int,
) ([]models.User, error) {
	users := make([]models.User, 0)
	err := r.db.
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&users).Error
	return users, err
}

func (r *UserRepository) FindAllCount() (int64, error) {
	var count int64 = 0
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	user := new(models.User)
	err := r.db.
		First(user, id).
		Error
	return user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
