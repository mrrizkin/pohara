package dialect

import (
	"fmt"
	"strings"
)

type SQLite struct{}

func (d SQLite) Placeholder() string { return "?" }
func (d SQLite) HasReturning() bool  { return false }
func (d SQLite) LimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return ""
}

func (d SQLite) QuoteIdentifier(name string) string {
	return fmt.Sprintf("`%s`", strings.ReplaceAll(name, "`", "``"))
}

func (d SQLite) ReplacePlaceholders(sql string) string {
	return sql // SQLite uses ? directly
}
