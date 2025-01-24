package provider

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/mysql/v2"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gofiber/storage/sqlite3/v2"

	"github.com/mrrizkin/pohara/app/config"
)

type Database struct {
	config *config.Config
}

func NewDatabase(config *config.Config) *Database {
	return &Database{config: config}
}

func (d *Database) Setup() (fiber.Storage, error) {
	switch d.config.Database.Driver {
	case "pgsql":
		return createPostgresStorage(d.config)
	case "mysql":
		return createMysqlStorage(d.config)
	case "sqlite":
		return createSQLiteStorage(d.config)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", d.config.Database.Driver)
	}
}

func createPostgresStorage(cfg *config.Config) (fiber.Storage, error) {
	config := postgres.Config{
		Host:       cfg.Database.Host,
		Port:       cfg.Database.Port,
		Database:   cfg.Database.Name,
		Username:   cfg.Database.Username,
		Password:   cfg.Database.Password,
		Table:      "sessions",
		SSLMode:    cfg.Database.SSLmode,
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return postgres.New(config), nil
}

func createMysqlStorage(cfg *config.Config) (fiber.Storage, error) {
	config := mysql.Config{
		Host:       cfg.Database.Host,
		Port:       cfg.Database.Port,
		Database:   cfg.Database.Name,
		Username:   cfg.Database.Username,
		Password:   cfg.Database.Password,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return mysql.New(config), nil
}

func createSQLiteStorage(cfg *config.Config) (fiber.Storage, error) {
	config := sqlite3.Config{
		Database:   cfg.Database.Host,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return sqlite3.New(config), nil
}
