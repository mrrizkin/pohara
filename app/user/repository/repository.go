package repository

import (
	"github.com/mrrizkin/pohara/app/user/entity"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"github.com/mrrizkin/pohara/internal/ports"
	"go.uber.org/fx"
)

type UserRepository struct {
	db ports.Database
}

type RepositoryDependencies struct {
	fx.In

	Db ports.Database
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

func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user)
}

func (r *UserRepository) FindAll(offset, limit sql.Int64Nullable) (*ports.FindAllResult, error) {
	var users []entity.User
	result, err := r.db.FindAll(&users, ports.Pagination{
		Limit:  limit,
		Offset: offset,
	})
	return result, err
}

func (r *UserRepository) FindAllCount() (int64, error) {
	return r.db.Count(&entity.User{})
}

func (r *UserRepository) FindByID(id uint) (*entity.User, error) {
	user := new(entity.User)
	err := r.db.First(user, id)
	return user, err
}

func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user)
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id)
}
