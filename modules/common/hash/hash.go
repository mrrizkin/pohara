package hash

import (
	"github.com/mrrizkin/pohara/modules/common/hash/algorithm"
)

type Algorithm interface {
	Generate(str string) (hash string, err error)
	Compare(password, hash string) (match bool, err error)
}

type Hashing struct {
	algorithm Algorithm
}

func New(config *Config) *Hashing {
	var algo Algorithm
	switch config.Provider {
	case "argon2", "argon":
		algo = algorithm.Argon2(
			uint32(config.Argon2.Memory),
			uint32(config.Argon2.Iterations),
			uint32(config.Argon2.KeyLen),
			uint32(config.Argon2.SaltLen),
			uint8(config.Argon2.Parallelism),
		)
	default:
		algo = algorithm.Argon2(
			uint32(config.Argon2.Memory),
			uint32(config.Argon2.Iterations),
			uint32(config.Argon2.KeyLen),
			uint32(config.Argon2.SaltLen),
			uint8(config.Argon2.Parallelism),
		)
	}

	return &Hashing{
		algorithm: algo,
	}
}

func (p *Hashing) Generate(str string) (hash string, err error) {
	return p.algorithm.Generate(str)
}

func (p *Hashing) Compare(password, hash string) (match bool, err error) {
	return p.algorithm.Compare(password, hash)
}
