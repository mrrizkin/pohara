package database

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/database/driver"
)

type Database struct {
	*gorm.DB
	config *config.Config
}

type Driver interface {
	Connect(cfg *config.Config) (*gorm.DB, error)
}

type DatabaseDependencies struct {
	fx.In

	Config *config.Config
}

func NewDatabase(lc fx.Lifecycle, deps DatabaseDependencies) *Database {
	db := &Database{}

	var d Driver
	switch deps.Config.Database.Driver {
	case "mysql", "mariadb", "maria":
		d = driver.Mysql{}
	case "pgsql", "postgres", "postgresql":
		d = driver.Postgres{}
	case "sqlite", "sqlite3", "file":
		d = driver.SQLite{}
	}
	gormDB, err := d.Connect(deps.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	db.DB = gormDB

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Close()
		},
	})

	return db
}

func (d *Database) Close() error {
	if d.DB == nil {
		return errors.New("you try to close database, but database not connected yet")
	}

	db, err := d.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
