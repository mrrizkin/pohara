package dialect

import (
	"fmt"
	"regexp"
)

type Dialect interface {
	Placeholder() string
	QuoteIdentifier(string) string
	ReplacePlaceholders(string) string
	HasReturning() bool
	LimitOffset(limit, offset int) string
}

func replacePlaceholders(sql, placeholder string) string {
	counter := 0
	return regexp.MustCompile(`\?`).ReplaceAllStringFunc(sql, func(string) string {
		counter++
		return fmt.Sprintf("%s%d", placeholder, counter)
	})
}
