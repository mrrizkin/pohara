package driver

import (
	"github.com/jmoiron/sqlx"
	_ "zombiezen.com/go/sqlite"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/core/logger"
)

type Sqlite struct {
	config *config.Config
	log    *logger.ZeroLog
}

func NewSqlite(
	config *config.Config,
	logger *logger.ZeroLog,
) *Sqlite {
	return &Sqlite{
		config: config,
		log:    logger,
	}
}

func (s *Sqlite) DSN() string {
	return s.config.Database.Host
}

func (s *Sqlite) Connect() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", s.DSN())
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		s.log.Warn("Failed to enable WAL journal mode", "error", err)
	} else {
		s.log.Info("Enabled WAL journal mode")
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		s.log.Warn("Failed to enable foreign keys", "error", err)
	} else {
		s.log.Info("Enabled foreign keys")
	}

	return db, nil
}
