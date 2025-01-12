package model

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/core/database"
)

var Module = fx.Module("model",
	fx.Invoke(func(db *database.GormDB) {
		db.AutoMigrate(
			&CfgPolicy{},
			&MRole{},
			&MUser{},
			&MUserAttribute{},
		)
	}),
)
