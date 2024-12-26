package hashing

import (
	"github.com/mrrizkin/pohara/config"

	"github.com/mrrizkin/pohara/internal/common/hashing/algorithm"
)

type Algorithm interface {
	Generate(str string) (hash string, err error)
	Compare(password, hash string) (match bool, err error)
}

type Hashing struct {
	algorithm Algorithm
}

func New(config *config.App) *Hashing {
	var algo Algorithm
	switch config.HASH_PROVIDER {
	case "argon2", "argon":
		algo = algorithm.Argon2(
			uint32(config.HASH_MEMORY),
			uint32(config.HASH_ITERATIONS),
			uint32(config.HASH_KEY_LEN),
			uint32(config.HASH_SALT_LEN),
			uint8(config.HASH_PARALLELISM),
		)
	default:
		algo = algorithm.Argon2(
			uint32(config.HASH_MEMORY),
			uint32(config.HASH_ITERATIONS),
			uint32(config.HASH_KEY_LEN),
			uint32(config.HASH_SALT_LEN),
			uint8(config.HASH_PARALLELISM),
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
