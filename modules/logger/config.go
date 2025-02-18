package logger

type Config struct {
	Name      string `env:"LOG_NAME,default=application.log"`
	Level     string `env:"LOG_LEVEL,default=debug"`
	Console   bool   `env:"LOG_CONSOLE,default=true"`
	File      bool   `env:"LOG_FILE,default=true"`
	Dir       string `env:"LOG_DIR"`
	MaxSize   int    `env:"LOG_MAX_SIZE,default=50"`
	MaxAge    int    `env:"LOG_MAX_AGE,default=7"`
	MaxBackup int    `env:"LOG_MAX_BACKUP,default=20"`
	Json      bool   `env:"LOG_JSON,default=true"`
}
