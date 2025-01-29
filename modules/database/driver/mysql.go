package driver

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/mrrizkin/pohara/app/config"
)

type Mysql struct {
	config *config.Config
}

func NewMysql(config *config.Config) *Mysql {
	return &Mysql{config: config}
}

func (m *Mysql) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.config.Database.Username,
		m.config.Database.Password,
		m.config.Database.Host,
		m.config.Database.Port,
		m.config.Database.Name,
	)
}

func (m *Mysql) Connect() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", m.DSN())
}
