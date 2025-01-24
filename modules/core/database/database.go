package database

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/core/database/provider"
	"github.com/mrrizkin/pohara/modules/core/logger"
)

type GormDB struct {
	*gorm.DB
}

type GormDBDriver interface {
	DSN() string
	Connect(cfg *config.Config) (*gorm.DB, error)
}

type GormDBDependencies struct {
	fx.In

	Config *config.Config
	Logger *logger.ZeroLog
}

func NewGormDB(
	lc fx.Lifecycle,
	deps GormDBDependencies,
) (*GormDB, error) {
	var driver GormDBDriver
	switch deps.Config.Database.Driver {
	case "mysql", "mariadb", "maria":
		driver = provider.NewMysql(deps.Config)
	case "pgsql", "postgres", "postgresql":
		driver = provider.NewPostgres(deps.Config)
	case "sqlite", "sqlite3", "file":
		driver = provider.NewSqlite(deps.Config, deps.Logger)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", deps.Config.Database.Driver)

	}

	db, err := driver.Connect(deps.Config)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}

			return sqlDB.Close()
		},
	})

	return &GormDB{db}, nil
}
