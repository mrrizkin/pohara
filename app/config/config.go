package config

import (
	"time"
)

type Config struct {
	App struct {
		Name    string `env:"APP_NAME,required"`
		Env     string `env:"APP_ENV,required"`
		URL     string `env:"APP_URL,required"`
		Port    int    `env:"APP_PORT,required"`
		Prefork bool   `env:"APP_PREFORK,default=false"`
		Debug   bool   `env:"APP_DEBUG,default=false"`
	}

	Database struct {
		Driver   string `env:"DB_DRIVER,default=sqlite"`
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT,default=5432"`
		Name     string `env:"DB_NAME"`
		Username string `env:"DB_USERNAME,default=root"`
		Password string `env:"DB_PASSWORD,default=root"`
		SSLmode  string `env:"DB_SSLMODE,default=disable"`
	}

	Site Site

	Log struct {
		Level     string `env:"LOG_LEVEL,default=debug"`
		Console   bool   `env:"LOG_CONSOLE,default=true"`
		File      bool   `env:"LOG_FILE,default=true"`
		Dir       string `env:"LOG_DIR"`
		MaxSize   int    `env:"LOG_MAX_SIZE,default=50"`
		MaxAge    int    `env:"LOG_MAX_AGE,default=7"`
		MaxBackup int    `env:"LOG_MAX_BACKUP,default=20"`
		Json      bool   `env:"LOG_JSON,default=true"`
	}

	Hash struct {
		Provider string `env:"HASH_PROVIDER,default=argon2"`

		Argon2 struct {
			Memory      int `env:"ARGON2_MEMORY,default=64"`
			Iterations  int `env:"ARGON2_ITERATIONS,default=10"`
			Parallelism int `env:"ARGON2_PARALLELISM,default=2"`
			SaltLen     int `env:"ARGON2_SALT_LEN,default=32"`
			KeyLen      int `env:"ARGON2_KEY_LEN,default=32"`
		}
	}

	CSRF struct {
		Key        string `env:"CSRF_KEY,default=X-CSRF-Token"`
		CookieName string `env:"CSRF_COOKIE_NAME,default=fiber_csrf_token"`
		SameSite   string `env:"CSRF_SAME_SITE,default=Lax"`
		Secure     bool   `env:"CSRF_SECURE,default=false"`
		Session    bool   `env:"CSRF_SESSION,default=true"`
		HttpOnly   bool   `env:"CSRF_HTTP_ONLY,default=true"`
		Expiration int    `env:"CSRF_EXPIRATION,default=3600"`
	}

	Session struct {
		Driver   string `env:"SESSION_DRIVER,default=file"`
		HttpOnly bool   `env:"SESSION_HTTP_ONLY,default=true"`
		Secure   bool   `env:"SESSION_SECURE,default=true"`
		SameSite string `env:"SESSION_SAME_SITE,default=Lax"`
	}

	View struct {
		Directory string `env:"VIEW_DIRECTORY,default=/views"`
		Extension string `env:"VIEW_EXTENSION,default=.html"`
		Cache     bool   `env:"VIEW_CACHE,default=true"`
	}

	Inertia struct {
		ManifestPath   string `env:"INERTIA_VITE_MANIFEST_PATH,default=public/build/manifest.json"`
		ContainerID    string `env:"INERTIA_CONTAINER_ID,default=app"`
		EncryptHistory bool   `env:"INERTIA_ENCRYPT_HISTORY,default=true"`
		EntryPath      string `env:"INERTIA_ENTRY_PATH,default=admin/index.html"`
	}

	StoragePath string `env:"STORAGE_PATH,default=storage"`
	SwaggerPath string `env:"SWAGGER_PATH,default=/docs/swagger.json"`
	CacheTTL    int    `env:"CACHE_TTL,default=300"`
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production" || c.App.Env == "prod"
}

func (c *Config) IsCacheView() bool {
	return c.View.Cache && c.IsProduction()
}

func (c *Config) CacheTTLSecond() time.Duration {
	return time.Duration(c.CacheTTL) * time.Second
}
