package driver

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mrrizkin/pohara/modules/database/config"
)

type Postgres struct{}

func (Postgres) Connect(config *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Name,
		config.Password,
		config.SSLmode,
	)))
}
