package hash

type Config struct {
	Provider string `env:"HASH_PROVIDER,default=argon2"`

	Argon2 struct {
		Memory      int `env:"ARGON2_MEMORY,default=64"`
		Iterations  int `env:"ARGON2_ITERATIONS,default=10"`
		Parallelism int `env:"ARGON2_PARALLELISM,default=2"`
		SaltLen     int `env:"ARGON2_SALT_LEN,default=32"`
		KeyLen      int `env:"ARGON2_KEY_LEN,default=32"`
	}
}
