package provider

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/config"

	"github.com/mrrizkin/pohara/internal/ports"
)

type Sqlite struct {
	config *config.Database
	log    ports.Logger
}

func NewSqlite(
	config *config.Database,
	log ports.Logger,
) *Sqlite {
	return &Sqlite{
		config: config,
		log:    log,
	}
}

func (s *Sqlite) DSN() string {
	return s.config.HOST
}

func (s *Sqlite) Connect(cfg *config.Database) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(s.DSN()))
	if err != nil {
		return nil, err
	}

	err = db.Exec("PRAGMA journal_mode = WAL;").Error
	if err != nil {
		s.log.Warn("Failed to enable WAL journal mode", "error", err)
	} else {
		s.log.Info("Enabled WAL journal mode")
	}

	err = db.Exec("PRAGMA foreign_keys = ON;").Error
	if err != nil {
		s.log.Warn("Failed to enable foreign keys", "error", err)
	} else {
		s.log.Info("Enabled foreign keys")
	}

	return db, nil
}
