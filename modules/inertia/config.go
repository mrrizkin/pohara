package inertia

type Config struct {
	ManifestPath   string `env:"INERTIA_VITE_MANIFEST_PATH,default=public/build/manifest.json"`
	ContainerID    string `env:"INERTIA_CONTAINER_ID,default=app"`
	EncryptHistory bool   `env:"INERTIA_ENCRYPT_HISTORY,default=true"`
	EntryPath      string `env:"INERTIA_ENTRY_PATH,default=admin/index.html"`
}
