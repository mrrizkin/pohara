package repository

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/cache"
	"github.com/mrrizkin/pohara/modules/database/db"
)

type AuthRepository struct {
	db    *db.Database
	cache *cache.Cache
}

type AuthRepositoryDependencies struct {
	fx.In

	DB    *db.Database
	Cache *cache.Cache
}

func NewAuthRepository(deps AuthRepositoryDependencies) *AuthRepository {
	return &AuthRepository{
		db:    deps.DB,
		cache: deps.Cache,
	}
}
