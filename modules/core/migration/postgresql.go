package migration

import (
	"fmt"
	"strings"
)

type PostgresDialect struct{}

func (d *PostgresDialect) GetStringType(length int) string {
	return fmt.Sprintf("VARCHAR(%d)", length)
}

func (d *PostgresDialect) GetIntegerType() string {
	return "INTEGER"
}

func (d *PostgresDialect) GetBigIntegerType() string {
	return "BIGINT"
}

func (d *PostgresDialect) GetTextType() string {
	return "TEXT"
}

func (d *PostgresDialect) GetTimestampType() string {
	return "TIMESTAMP"
}

func (d *PostgresDialect) GetBooleanType() string {
	return "BOOLEAN"
}

func (d *PostgresDialect) GetEnumType(values []string) string {
	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("TEXT CHECK (%s IN (%s))", "value", strings.Join(quotedValues, ", "))
}

func (d *PostgresDialect) GetDecimalType(precision, scale int) string {
	return fmt.Sprintf("DECIMAL(%d, %d)", precision, scale)
}

func (d *PostgresDialect) GetDateType() string {
	return "DATE"
}

func (d *PostgresDialect) GetTimeType() string {
	return "TIME"
}

func (d *PostgresDialect) GetJsonType() string {
	return "JSONB"
}

func (d *PostgresDialect) GetBinaryType() string {
	return "BYTEA"
}

func (d *PostgresDialect) GetFloatType() string {
	return "FLOAT"
}

func (d *PostgresDialect) GetUUIDType() string {
	return "UUID"
}

func (d *PostgresDialect) GetPrimaryKey() string {
	return "PRIMARY KEY"
}

func (d *PostgresDialect) GetUnique(column string) string {
	return "UNIQUE"
}

func (d *PostgresDialect) GetNullable(column string, nullable bool) string {
	if nullable {
		return "NULL"
	}
	return "NOT NULL"
}

func (d *PostgresDialect) GetDefault(column, value string) string {
	return fmt.Sprintf("DEFAULT '%s'", value)
}

func (d *PostgresDialect) CreateTableSQL(table string, columns []string) string {
	return fmt.Sprintf("CREATE TABLE %s (\n  %s\n)", table, strings.Join(columns, ",\n  "))
}

func (d *PostgresDialect) AddColumnSQL(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", table, column)
}

func (d *PostgresDialect) ModifyColumnSQL(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s", table, column)
}

func (d *PostgresDialect) RenameColumnSQL(table, old, new string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", table, old, new)
}

func (d *PostgresDialect) DropColumnSQL(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, column)
}

func (d *PostgresDialect) DropTableSQL(table string) string {
	return fmt.Sprintf("DROP TABLE %s", table)
}

func (d *PostgresDialect) RenameTableSQL(old, new string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s", old, new)
}

func (d *PostgresDialect) CreateIndexSQL(table, name string, columns []string, unique bool) string {
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE"
	}
	return fmt.Sprintf(
		"CREATE %s INDEX %s ON %s (%s)",
		uniqueStr,
		name,
		table,
		strings.Join(columns, ", "),
	)
}

func (d *PostgresDialect) DropIndexSQL(table, name string) string {
	return fmt.Sprintf("DROP INDEX %s", name)
}
