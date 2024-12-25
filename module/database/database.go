package database

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/module/database/provider"
	"github.com/mrrizkin/pohara/module/logger"
)

type DatabaseDriver interface {
	DSN() string
	Connect(cfg *config.Database) (*gorm.DB, error)
}

type Database struct {
	*gorm.DB
}

type Dependencies struct {
	fx.In

	ConfigDB *config.Database
	Log      *logger.Logger
}

type Result struct {
	fx.Out

	Database *Database
}

func New(
	lc fx.Lifecycle,
	deps Dependencies,
) (Result, error) {
	var driver DatabaseDriver
	switch deps.ConfigDB.DRIVER {
	case "mysql", "mariadb", "maria":
		driver = provider.NewMysql(deps.ConfigDB)
	case "pgsql", "postgres", "postgresql":
		driver = provider.NewPostgres(deps.ConfigDB)
	case "sqlite", "sqlite3", "file":
		driver = provider.NewSqlite(deps.ConfigDB, deps.Log)
	default:
		return Result{}, fmt.Errorf("unknown database driver: %s", deps.ConfigDB.DRIVER)

	}

	db, err := driver.Connect(deps.ConfigDB)
	if err != nil {
		return Result{}, err
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

	return Result{
		Database: &Database{DB: db},
	}, nil
}

func (d *Database) Stop() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
