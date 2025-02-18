package config

type Config struct {
	CookieName string `env:"SESSION_COOKIE_NAME,default=pohara_session_key"`
	Driver     string `env:"SESSION_DRIVER,default=file"`
	HttpOnly   bool   `env:"SESSION_HTTP_ONLY,default=true"`
	Secure     bool   `env:"SESSION_SECURE,default=true"`
	SameSite   string `env:"SESSION_SAME_SITE,default=Lax"`
}
