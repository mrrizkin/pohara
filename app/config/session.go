package config

type Session struct {
	DRIVER    string `env:"SESSION_DRIVER,default=file"`
	HTTP_ONLY bool   `env:"SESSION_HTTP_ONLY,default=true"`
	SECURE    bool   `env:"SESSION_SECURE,default=true"`
	SAME_SITE string `env:"SESSION_SAME_SITE,default=Lax"`
}
