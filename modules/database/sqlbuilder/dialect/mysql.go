package dialect

import (
	"fmt"
	"strings"
)

type MySQL struct{}

func (d MySQL) Placeholder() string { return "?" }
func (d MySQL) HasReturning() bool  { return false }
func (d MySQL) LimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return ""
}

func (d MySQL) QuoteIdentifier(name string) string {
	return fmt.Sprintf("`%s`", strings.ReplaceAll(name, "`", "``"))
}

func (d MySQL) ReplacePlaceholders(sql string) string {
	return sql // MySQL uses ? directly
}
