package driver

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/app/config"
)

type SQLite struct{}

func (SQLite) Connect(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Database.Host))
	if err != nil {
		return nil, err
	}

	if err := db.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
		return nil, err
	}

	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return nil, err
	}

	return db, nil
}
