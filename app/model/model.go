package model

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module("model",
	fx.Invoke(func(db *gorm.DB) {
		db.AutoMigrate(
			&User{},
		)
	}),
)
