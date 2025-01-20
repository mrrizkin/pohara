package migration

import (
	"fmt"
	"strings"
)

// PostgresDialect implements Dialect for PostgreSQL
type PostgresDialect struct{}

func (p *PostgresDialect) GetStringType(length int) string {
	return fmt.Sprintf("VARCHAR(%d)", length)
}

func (p *PostgresDialect) GetIntegerType() string {
	return "INTEGER"
}

func (p *PostgresDialect) GetBigIntegerType() string {
	return "BIGINT"
}

func (p *PostgresDialect) GetTextType() string {
	return "TEXT"
}

func (p *PostgresDialect) GetTimestampType() string {
	return "TIMESTAMP"
}

func (p *PostgresDialect) GetBooleanType() string {
	return "BOOLEAN"
}

func (p *PostgresDialect) GetDecimalType(precision, scale int) string {
	return fmt.Sprintf("DECIMAL(%d,%d)", precision, scale)
}

func (p *PostgresDialect) GetJsonType() string {
	return "JSONB"
}

func (p *PostgresDialect) GetDateType() string {
	return "DATE"
}

func (p *PostgresDialect) GetTimeType() string {
	return "TIME"
}

func (p *PostgresDialect) GetBinaryType() string {
	return "BYTEA"
}

func (p *PostgresDialect) GetFloatType() string {
	return "FLOAT"
}

func (p *PostgresDialect) GetUUIDType() string {
	return "UUID"
}

func (p *PostgresDialect) WrapPrimaryKey(column string) string {
	return "PRIMARY KEY"
}

func (p *PostgresDialect) WrapUnique(column string) string {
	return "UNIQUE"
}

func (p *PostgresDialect) WrapNullable(column string) string {
	return "NOT NULL"
}

func (p *PostgresDialect) WrapDefault(column, value string) string {
	return fmt.Sprintf("DEFAULT %s", value)
}

func (p *PostgresDialect) GetForeignKeyConstraint(column, refTable, refColumn, onDelete, onUpdate string) string {
	constraint := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", column, refTable, refColumn)
	if onDelete != "" {
		constraint += " ON DELETE " + onDelete
	}
	if onUpdate != "" {
		constraint += " ON UPDATE " + onUpdate
	}
	return constraint
}

func (p *PostgresDialect) CreateIndex(table, name string, columns []string, unique bool) string {
	if name == "" {
		name = fmt.Sprintf("%s_%s_idx", table, strings.Join(columns, "_"))
	}
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}
	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)", uniqueStr, name, table, strings.Join(columns, ", "))
}
