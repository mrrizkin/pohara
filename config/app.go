package config

import "time"

type App struct {
	APP_NAME    string `env:"APP_NAME,required"`
	APP_KEY     string `env:"APP_KEY"`
	APP_ENV     string `env:"ENV,required"`
	APP_URL     string `env:"URL,required"`
	APP_PORT    int    `env:"PORT,required"`
	APP_PREFORK bool   `env:"PREFORK,default=false"`
	APP_DEBUG   bool   `env:"DEBUG,default=false"`

	LOG_LEVEL      string `env:"LOG_LEVEL,default=debug"`
	LOG_CONSOLE    bool   `env:"LOG_CONSOLE,default=true"`
	LOG_FILE       bool   `env:"LOG_FILE,default=true"`
	LOG_DIR        string `env:"LOG_DIR"`
	LOG_MAX_SIZE   int    `env:"LOG_MAX_SIZE,default=50"`
	LOG_MAX_AGE    int    `env:"LOG_MAX_AGE,default=7"`
	LOG_MAX_BACKUP int    `env:"LOG_MAX_BACKUP,default=20"`
	LOG_JSON       bool   `env:"LOG_JSON,default=true"`

	HASH_PROVIDER    string `env:"HASH_PROVIDER,default=argon2"`
	HASH_MEMORY      int    `env:"HASH_MEMORY,default=64"`
	HASH_ITERATIONS  int    `env:"HASH_ITERATIONS,default=10"`
	HASH_PARALLELISM int    `env:"HASH_PARALLELISM,default=2"`
	HASH_SALT_LEN    int    `env:"HASH_SALT_LEN,default=32"`
	HASH_KEY_LEN     int    `env:"HASH_KEY_LEN,default=32"`

	CSRF_KEY         string `env:"CSRF_KEY,default=X-CSRF-Token"`
	CSRF_COOKIE_NAME string `env:"CSRF_COOKIE_NAME,default=fiber_csrf_token"`
	CSRF_SAME_SITE   string `env:"CSRF_SAME_SITE,default=Lax"`
	CSRF_SECURE      bool   `env:"CSRF_SECURE,default=false"`
	CSRF_SESSION     bool   `env:"CSRF_SESSION,default=true"`
	CSRF_HTTP_ONLY   bool   `env:"CSRF_HTTP_ONLY,default=true"`
	CSRF_EXPIRATION  int    `env:"CSRF_EXPIRATION,default=3600"`

	VIEW_DIRECTORY string `env:"VIEW_DIRECTORY,default=/views"`
	VIEW_EXTENSION string `env:"VIEW_EXTENSION,default=.html"`
	VIEW_CACHE     bool   `env:"VIEW_CACHE,default=true"`

	STORAGE_PATH string `env:"STORAGE_PATH,default=storage"`
	SWAGGER_PATH string `env:"SWAGGER_PATH,default=/docs/swagger.json"`
	CACHE_TTL    int    `env:"CACHE_TTL,default=300"`
}

func (a *App) IsProduction() bool {
	return a.APP_ENV == "production" || a.APP_ENV == "prod"
}

func (a *App) CacheTTLSecond() time.Duration {
	return time.Duration(a.CACHE_TTL) * time.Second
}
