package config

type Config struct {
	Driver   string `env:"DB_DRIVER,default=sqlite"`
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT,default=5432"`
	Name     string `env:"DB_NAME"`
	Username string `env:"DB_USERNAME,default=root"`
	Password string `env:"DB_PASSWORD,default=root"`
	SSLmode  string `env:"DB_SSLMODE,default=disable"`
}
