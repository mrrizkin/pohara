package hash

import gonanoid "github.com/matoous/go-nanoid/v2"

func NanoID() string {
	id, _ := gonanoid.New()
	return id
}

func GenerateNanoID(chars string, size int) string {
	id, _ := gonanoid.Generate(chars, size)
	return id
}
