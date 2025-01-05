package nanoid

import gonanoid "github.com/matoous/go-nanoid/v2"

func New() string {
	id, _ := gonanoid.New()
	return id
}

func Generate(chars string, size int) string {
	id, _ := gonanoid.Generate(chars, size)
	return id
}
