package models

import (
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/module/database"
	"github.com/mrrizkin/pohara/module/logger"
	"go.uber.org/fx"
)

type Model struct {
	db     *database.Database
	config *config.Database
	log    *logger.Logger

	models []interface{}
}

type Dependencies struct {
	fx.In

	Db     *database.Database
	Config *config.Database
	Log    *logger.Logger
}

type Result struct {
	fx.Out

	Model *Model
}

func New(deps Dependencies) Result {
	return Result{
		Model: &Model{
			db:     deps.Db,
			config: deps.Config,
			log:    deps.Log,

			models: []interface{}{
				&User{},
			},
		},
	}
}

var Module = fx.Module("models",
	fx.Provide(New),
	fx.Invoke(func(model *Model) {
		model.Migrate()
	}),
)

func (m *Model) Migrate() error {
	if len(m.models) == 0 {
		return nil
	}

	if !m.config.AUTO_MIGRATE {
		return nil
	}

	m.log.Info("migrating models", "count", len(m.models))
	return m.db.AutoMigrate(m.models...)
}

func (m *Model) Seed() error {
	return nil
}
