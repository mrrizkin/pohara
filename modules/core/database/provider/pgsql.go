package provider

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mrrizkin/pohara/app/config"
)

type Postgres struct {
	config *config.Config
}

func NewPostgres(config *config.Config) *Postgres {
	return &Postgres{config: config}
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		p.config.Database.Host,
		p.config.Database.Port,
		p.config.Database.Username,
		p.config.Database.Name,
		p.config.Database.Password,
		p.config.Database.SSLmode,
	)
}

func (p *Postgres) Connect(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(p.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
