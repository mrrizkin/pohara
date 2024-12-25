package user

import (
	"fmt"

	"github.com/mrrizkin/pohara/models"
	"github.com/mrrizkin/pohara/module/database/sql"
	"github.com/mrrizkin/pohara/module/hashing"
	"go.uber.org/fx"
)

type UserService struct {
	repo    *UserRepository
	hashing *hashing.Hashing
}

type ServiceDependencies struct {
	fx.In

	UserRepository *UserRepository
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

func (s *UserService) Create(user *models.User) (*models.User, error) {
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

func (s *UserService) FindAll(page, perPage int) (map[string]interface{}, error) {
	users, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	usersCount, err := s.repo.FindAllCount()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"result": users,
		"total":  int(usersCount),
	}, nil
}

func (s *UserService) FindByID(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(id uint, user *models.User) (*models.User, error) {
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
