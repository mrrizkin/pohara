package driver

import (
	"fmt"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"

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

func (p *Postgres) Connect() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", p.DSN())
}
