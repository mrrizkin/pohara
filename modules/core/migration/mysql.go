package migration

import (
	"fmt"
	"strings"
)

// MySQL dialect implementation
type MySqlDialect struct{}

func (d *MySqlDialect) GetStringType(length int) string {
	return fmt.Sprintf("VARCHAR(%d)", length)
}

func (d *MySqlDialect) GetIntegerType() string {
	return "INT"
}

func (d *MySqlDialect) GetBigIntegerType() string {
	return "BIGINT UNSIGNED"
}

func (d *MySqlDialect) GetTextType() string {
	return "TEXT"
}

func (d *MySqlDialect) GetTimestampType() string {
	return "TIMESTAMP"
}

func (d *MySqlDialect) GetBooleanType() string {
	return "BOOLEAN"
}

func (d *MySqlDialect) GetEnumType(values []string) string {
	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("ENUM(%s)", strings.Join(quotedValues, ", "))
}

func (d *MySqlDialect) GetDecimalType(precision, scale int) string {
	return fmt.Sprintf("DECIMAL(%d, %d)", precision, scale)
}

func (d *MySqlDialect) GetDateType() string {
	return "DATE"
}

func (d *MySqlDialect) GetTimeType() string {
	return "TIME"
}

func (d *MySqlDialect) GetJsonType() string {
	return "JSON"
}

func (d *MySqlDialect) GetBinaryType() string {
	return "BLOB"
}

func (d *MySqlDialect) GetFloatType() string {
	return "FLOAT"
}

func (d *MySqlDialect) GetUUIDType() string {
	return "CHAR(36)"
}

func (d *MySqlDialect) GetPrimaryKey() string {
	return "AUTO_INCREMENT PRIMARY KEY"
}

func (d *MySqlDialect) GetUnique(column string) string {
	return "UNIQUE"
}

func (d *MySqlDialect) GetNullable(column string, nullable bool) string {
	if nullable {
		return "NULL"
	}
	return "NOT NULL"
}

func (d *MySqlDialect) GetDefault(column, value string) string {
	return fmt.Sprintf("DEFAULT '%s'", value)
}

func (d *MySqlDialect) CreateTableSQL(table string, columns []string) string {
	return fmt.Sprintf("CREATE TABLE %s (\n  %s\n)", table, strings.Join(columns, ",\n  "))
}

func (d *MySqlDialect) CreateTableIfNotExistSQL(table string, columns []string) string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXIST %s (\n  %s\n)",
		table,
		strings.Join(columns, ",\n  "),
	)
}

func (d *MySqlDialect) AddColumnSQL(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", table, column)
}

func (d *MySqlDialect) ModifyColumnSQL(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s", table, column)
}

func (d *MySqlDialect) RenameColumnSQL(table, old, new string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", table, old, new)
}

func (d *MySqlDialect) DropColumnSQL(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, column)
}

func (d *MySqlDialect) DropTableSQL(table string) string {
	return fmt.Sprintf("DROP TABLE %s", table)
}

func (d *MySqlDialect) RenameTableSQL(old, new string) string {
	return fmt.Sprintf("RENAME TABLE %s TO %s", old, new)
}

func (d *MySqlDialect) CreateIndexSQL(table, name string, columns []string, unique bool) string {
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

func (d *MySqlDialect) DropIndexSQL(table, name string) string {
	return fmt.Sprintf("DROP INDEX %s ON %s", name, table)
}
