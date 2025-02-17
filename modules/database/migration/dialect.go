package migration

// Dialect interface implementation
type Dialect interface {
	// Column types
	GetStringType(length int) string
	GetIntegerType() string
	GetBigIntegerType() string
	GetTextType() string
	GetTimestampType() string
	GetBooleanType() string
	GetEnumType(values []string) string
	GetDecimalType(precision, scale int) string
	GetDateType() string
	GetTimeType() string
	GetJsonType() string
	GetBinaryType() string
	GetFloatType() string
	GetUUIDType() string

	// Constraints
	GetPrimaryKey() string
	GetUnique(column string) string
	GetNullable(column string, nullable bool) string
	GetDefault(column, value string) string

	// Table operations
	CreateTableSQL(table string, columns []string) string
	CreateTableIfNotExistSQL(table string, columns []string) string
	AddColumnSQL(table, column string) string
	ModifyColumnSQL(table, column string) string
	RenameColumnSQL(table, old, new string) string
	DropColumnSQL(table, column string) string
	DropTableSQL(table string) string
	DropTableIfExistSQL(table string) string
	RenameTableSQL(old, new string) string
	CreateIndexSQL(table, name string, columns []string, unique bool) string
	DropIndexSQL(table, name string) string
	CreateCompositePrimarySQL(columns ...string) string
}
