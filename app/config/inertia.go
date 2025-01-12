package config

type Inertia struct {
	MANIFEST_PATH   string `env:"INERTIA_VITE_MANIFEST_PATH,default=public/build/manifest.json"`
	CONTAINER_ID    string `env:"INERTIA_CONTAINER_ID,default=app"`
	ENCRYPT_HISTORY bool   `env:"INERTIA_ENCRYPT_HISTORY,default=true"`
	ENTRY_PATH      string `env:"INERTIA_ENTRY_PATH,default=index.html"`
}
