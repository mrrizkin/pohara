package migration

import (
	"fmt"
	"strings"
)

type SQLiteDialect struct{}

func (d *SQLiteDialect) GetStringType(length int) string {
	return "TEXT"
}

func (d *SQLiteDialect) GetIntegerType() string {
	return "INTEGER"
}

func (d *SQLiteDialect) GetBigIntegerType() string {
	return "INTEGER"
}

func (d *SQLiteDialect) GetTextType() string {
	return "TEXT"
}

func (d *SQLiteDialect) GetTimestampType() string {
	return "DATETIME"
}

func (d *SQLiteDialect) GetBooleanType() string {
	return "BOOLEAN"
}

func (d *SQLiteDialect) GetEnumType(values []string) string {
	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("TEXT CHECK (%s IN (%s))", "value", strings.Join(quotedValues, ", "))
}

func (d *SQLiteDialect) GetDecimalType(precision, scale int) string {
	return "REAL"
}

func (d *SQLiteDialect) GetDateType() string {
	return "DATE"
}

func (d *SQLiteDialect) GetTimeType() string {
	return "TIME"
}

func (d *SQLiteDialect) GetJsonType() string {
	return "TEXT"
}

func (d *SQLiteDialect) GetBinaryType() string {
	return "BLOB"
}

func (d *SQLiteDialect) GetFloatType() string {
	return "REAL"
}

func (d *SQLiteDialect) GetUUIDType() string {
	return "TEXT"
}

func (d *SQLiteDialect) GetPrimaryKey() string {
	return "PRIMARY KEY AUTOINCREMENT"
}

func (d *SQLiteDialect) GetUnique(column string) string {
	return "UNIQUE"
}

func (d *SQLiteDialect) GetNullable(column string, nullable bool) string {
	if nullable {
		return "NULL"
	}
	return "NOT NULL"
}

func (d *SQLiteDialect) GetDefault(column, value string) string {
	return fmt.Sprintf("DEFAULT '%s'", value)
}

func (d *SQLiteDialect) CreateTableSQL(table string, columns []string) string {
	return fmt.Sprintf("CREATE TABLE %s (\n  %s\n)", table, strings.Join(columns, ",\n  "))
}

func (d *SQLiteDialect) AddColumnSQL(table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", table, column)
}

func (d *SQLiteDialect) ModifyColumnSQL(table, column string) string {
	// SQLite does not support MODIFY COLUMN directly. You need to recreate the table.
	return fmt.Sprintf("ALTER TABLE %s RENAME TO old_table", table)
}

func (d *SQLiteDialect) RenameColumnSQL(table, old, new string) string {
	// SQLite does not support RENAME COLUMN directly. You need to recreate the table.
	return fmt.Sprintf("ALTER TABLE %s RENAME TO old_table", table)
}

func (d *SQLiteDialect) DropColumnSQL(table, column string) string {
	// SQLite does not support DROP COLUMN directly. You need to recreate the table.
	return fmt.Sprintf("ALTER TABLE %s RENAME TO old_table", table)
}

func (d *SQLiteDialect) DropTableSQL(table string) string {
	return fmt.Sprintf("DROP TABLE %s", table)
}

func (d *SQLiteDialect) RenameTableSQL(old, new string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s", old, new)
}

func (d *SQLiteDialect) CreateIndexSQL(table, name string, columns []string, unique bool) string {
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

func (d *SQLiteDialect) DropIndexSQL(table, name string) string {
	return fmt.Sprintf("DROP INDEX %s", name)
}
