package repository

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/core/cache"
	"github.com/mrrizkin/pohara/modules/database"
)

type AuthRepository struct {
	db    *database.Database
	cache *cache.Cache
}

type AuthRepositoryDependencies struct {
	fx.In

	DB    *database.Database
	Cache *cache.Cache
}

func NewAuthRepository(deps AuthRepositoryDependencies) *AuthRepository {
	return &AuthRepository{
		db:    deps.DB,
		cache: deps.Cache,
	}
}
