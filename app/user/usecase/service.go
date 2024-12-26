package usecase

import (
	"fmt"

	"github.com/mrrizkin/pohara/app/user/entity"
	"github.com/mrrizkin/pohara/app/user/repository"
	"github.com/mrrizkin/pohara/internal/common/hashing"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"github.com/mrrizkin/pohara/internal/ports"
	"go.uber.org/fx"
)

type UserService struct {
	repo    *repository.UserRepository
	hashing *hashing.Hashing
}

type ServiceDependencies struct {
	fx.In

	UserRepository *repository.UserRepository
	Hashing        *hashing.Hashing
}

type ServiceResult struct {
	fx.Out

	UserService *UserService
}

func Service(deps ServiceDependencies) ServiceResult {
	return ServiceResult{
		UserService: &UserService{
			repo:    deps.UserRepository,
			hashing: deps.Hashing,
		},
	}
}

func (s *UserService) Create(user *entity.User) (*entity.User, error) {
	if !user.Password.Valid {
		return nil, fmt.Errorf("password is required")
	}

	hash, err := s.hashing.Generate(user.Password.String)
	if err != nil {
		return nil, err
	}

	user.Password = sql.StringNullable{
		String: hash,
		Valid:  true,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindAll(offset, limit sql.Int64Nullable) (*ports.FindAllResult, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *UserService) FindByID(id uint) (*entity.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(id uint, user *entity.User) (*entity.User, error) {
	userExist, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user.Password.Valid {
		hash, err := s.hashing.Generate(user.Password.String)
		if err != nil {
			return nil, err
		}

		userExist.Password = sql.StringNullable{
			Valid:  true,
			String: hash,
		}
	}

	userExist.Name = user.Name
	userExist.Email = user.Email
	userExist.Username = user.Username

	err = s.repo.Update(userExist)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
