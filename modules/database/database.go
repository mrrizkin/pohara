package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/core/logger"
	"github.com/mrrizkin/pohara/modules/database/driver"
	"github.com/mrrizkin/pohara/modules/database/sqlbuilder/builder"
	"github.com/mrrizkin/pohara/modules/database/sqlbuilder/dialect"
)

type Database struct {
	*sqlx.DB
	config *config.Config
}

type DatabaseDriver interface {
	DSN() string
	Connect() (*sqlx.DB, error)
}

type DatabaseDependencies struct {
	fx.In

	Config *config.Config
	Logger *logger.ZeroLog
}

func NewDatabase(lc fx.Lifecycle, deps DatabaseDependencies) (*Database, error) {
	var databaseDriver DatabaseDriver

	switch deps.Config.Database.Driver {
	case "mysql", "mariadb", "maria":
		databaseDriver = driver.NewMysql(deps.Config)
	case "pgsql", "postgres", "postgresql":
		databaseDriver = driver.NewPostgres(deps.Config)
	case "sqlite", "sqlite3", "file":
		databaseDriver = driver.NewSqlite(deps.Config, deps.Logger)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", deps.Config.Database.Driver)
	}

	db, err := databaseDriver.Connect()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Close()
		},
	})

	return &Database{
		DB: db,
	}, err
}

func (d *Database) Builder() *builder.SQLBuilder {
	var sqlBuilder *builder.SQLBuilder

	switch d.config.Database.Driver {
	case "mysql", "mariadb", "maria":
		sqlBuilder = builder.New(dialect.MySQL{}, d.DB)
	case "pgsql", "postgres", "postgresql":
		sqlBuilder = builder.New(dialect.Postgres{}, d.DB)
	case "sqlite", "sqlite3", "file":
		sqlBuilder = builder.New(dialect.SQLite{}, d.DB)
	default:
		panic(fmt.Sprintf("unknown database driver: %s", d.config.Database.Driver))
	}

	return sqlBuilder
}
