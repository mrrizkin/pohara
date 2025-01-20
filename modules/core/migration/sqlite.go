package migration

import (
	"fmt"
	"strings"
)

// SQLiteDialect implements Dialect for SQLite
type SQLiteDialect struct{}

func (s *SQLiteDialect) GetStringType(length int) string {
	return "TEXT"
}

func (s *SQLiteDialect) GetIntegerType() string {
	return "INTEGER"
}

func (s *SQLiteDialect) GetBigIntegerType() string {
	return "INTEGER"
}

func (s *SQLiteDialect) GetTextType() string {
	return "TEXT"
}

func (s *SQLiteDialect) GetTimestampType() string {
	return "DATETIME"
}

func (s *SQLiteDialect) GetBooleanType() string {
	return "INTEGER"
}

func (s *SQLiteDialect) GetDecimalType(precision, scale int) string {
	return "REAL"
}

func (s *SQLiteDialect) GetJsonType() string {
	return "TEXT"
}

func (s *SQLiteDialect) GetDateType() string {
	return "DATE"
}

func (s *SQLiteDialect) GetTimeType() string {
	return "TIME"
}

func (s *SQLiteDialect) GetBinaryType() string {
	return "BLOB"
}

func (s *SQLiteDialect) GetFloatType() string {
	return "REAL"
}

func (s *SQLiteDialect) GetUUIDType() string {
	return "TEXT"
}

func (s *SQLiteDialect) WrapPrimaryKey(column string) string {
	return "PRIMARY KEY AUTOINCREMENT"
}

func (s *SQLiteDialect) WrapUnique(column string) string {
	return "UNIQUE"
}

func (s *SQLiteDialect) WrapNullable(column string) string {
	return "NOT NULL"
}

func (s *SQLiteDialect) WrapDefault(column, value string) string {
	switch strings.ToLower(value) {
	case "true":
		return "DEFAULT 1"
	case "false":
		return "DEFAULT 0"
	default:
		return fmt.Sprintf("DEFAULT %s", value)
	}
}

func (s *SQLiteDialect) GetForeignKeyConstraint(column, refTable, refColumn, onDelete, onUpdate string) string {
	constraint := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", column, refTable, refColumn)
	if onDelete != "" {
		constraint += " ON DELETE " + onDelete
	}
	if onUpdate != "" {
		constraint += " ON UPDATE " + onUpdate
	}
	return constraint
}

func (s *SQLiteDialect) CreateIndex(table, name string, columns []string, unique bool) string {
	if name == "" {
		name = fmt.Sprintf("%s_%s_idx", table, strings.Join(columns, "_"))
	}
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}
	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)",
		uniqueStr, name, table, strings.Join(columns, ", "))
}
