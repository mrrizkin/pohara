package dialect

import (
	"fmt"
	"strings"
)

type Postgres struct{}

func (d Postgres) Placeholder() string { return "$" }
func (d Postgres) HasReturning() bool  { return true }
func (d Postgres) LimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}

func (d Postgres) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(name, `"`, `""`))
}

func (d Postgres) ReplacePlaceholders(sql string) string {
	return replacePlaceholders(sql, d.Placeholder())
}
