package migration

import (
	"fmt"
	"strings"
)

// Dialect defines the interface for database-specific SQL generation
type Dialect interface {
	// Column type methods
	GetStringType(length int) string
	GetIntegerType() string
	GetBigIntegerType() string
	GetTextType() string
	GetTimestampType() string
	GetBooleanType() string
	GetDecimalType(precision, scale int) string
	GetJsonType() string
	GetDateType() string
	GetTimeType() string
	GetBinaryType() string
	GetFloatType() string
	GetUUIDType() string

	// Constraint methods
	WrapPrimaryKey(column string) string
	WrapUnique(column string) string
	WrapNullable(column string) string
	WrapDefault(column, value string) string

	// Foreign key methods
	GetForeignKeyConstraint(column, refTable, refColumn, onDelete, onUpdate string) string

	// Index methods
	CreateIndex(table, name string, columns []string, unique bool) string
}

// Column represents a database column definition
type Column struct {
	NameField       string
	TypeField       string
	NullableField   bool
	PrimaryField    bool
	UniqueField     bool
	DefaultField    *string
	ReferencesField *Reference
}

// Reference represents a foreign key reference
type Reference struct {
	TableField    string
	ColumnField   string
	OnDeleteField string
	OnUpdateField string
}

// Index represents a table index
type Index struct {
	NameField    string
	ColumnsField []string
	UniqueField  bool
}

// Schema represents the migration schema builder
type Schema struct {
	dialectData Dialect
	createSql   []string
}

// Blueprint represents a table schema definition
type Blueprint struct {
	tableName    string
	columnsData  []Column
	indexesData  []Index
	commandsData []string
	dialectData  Dialect
}

// NewSchema creates a new schema builder with the specified dialect
func NewSchema(dialect Dialect) *Schema {
	return &Schema{
		dialectData: dialect,
	}
}

// Create creates a new blueprint for table creation
func (s *Schema) Create(table string, fn func(*Blueprint)) *Schema {
	blueprint := &Blueprint{
		tableName:   table,
		dialectData: s.dialectData,
	}
	fn(blueprint)
	s.createSql = append(s.createSql, blueprint.ToSQL())
	return s
}

// Table modifies an existing table
func (s *Schema) Table(table string, fn func(*Blueprint)) string {
	blueprint := &Blueprint{
		tableName:   table,
		dialectData: s.dialectData,
	}
	fn(blueprint)
	return blueprint.ToSQL()
}

// Column methods
func (b *Blueprint) String(name string, length ...int) *Column {
	l := 255
	if len(length) > 0 {
		l = length[0]
	}
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetStringType(l),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Integer(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetIntegerType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) BigInteger(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetBigIntegerType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Text(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetTextType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Timestamp(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetTimestampType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Boolean(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetBooleanType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Decimal(name string, precision, scale int) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetDecimalType(precision, scale),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Float(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetFloatType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Json(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetJsonType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Date(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetDateType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Time(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetTimeType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Binary(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetBinaryType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) UUID(name string) *Column {
	col := &Column{
		NameField: name,
		TypeField: b.dialectData.GetUUIDType(),
	}
	b.columnsData = append(b.columnsData, *col)
	return col
}

func (b *Blueprint) Timestamps() {
	b.Timestamp("created_at").Nullable()
	b.Timestamp("updated_at").Nullable()
}

func (b *Blueprint) SoftDeletes() {
	b.Timestamp("deleted_at").Nullable()
}

// Column modifier methods
func (c *Column) Nullable() *Column {
	c.NullableField = true
	return c
}

func (c *Column) Primary() *Column {
	c.PrimaryField = true
	return c
}

func (c *Column) Unique() *Column {
	c.UniqueField = true
	return c
}

func (c *Column) Default(value string) *Column {
	c.DefaultField = &value
	return c
}

func (c *Column) References(table string) *Column {
	c.ReferencesField = &Reference{
		TableField:  table,
		ColumnField: "id",
	}
	return c
}

func (c *Column) On(column string) *Column {
	if c.ReferencesField != nil {
		c.ReferencesField.ColumnField = column
	}
	return c
}

func (c *Column) OnDelete(action string) *Column {
	if c.ReferencesField != nil {
		c.ReferencesField.OnDeleteField = action
	}
	return c
}

func (c *Column) OnUpdate(action string) *Column {
	if c.ReferencesField != nil {
		c.ReferencesField.OnUpdateField = action
	}
	return c
}

// Index methods
func (b *Blueprint) Index(columns ...string) {
	b.indexesData = append(b.indexesData, Index{
		ColumnsField: columns,
		UniqueField:  false,
	})
}

func (b *Blueprint) Unique(columns ...string) {
	b.indexesData = append(b.indexesData, Index{
		ColumnsField: columns,
		UniqueField:  true,
	})
}

// Drop and rename methods
func (b *Blueprint) DropColumn(name string) {
	b.commandsData = append(b.commandsData, fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", b.tableName, name))
}

func (b *Blueprint) RenameColumn(from, to string) {
	b.commandsData = append(b.commandsData, fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", b.tableName, from, to))
}

func (b *Blueprint) DropIndex(name string) {
	b.commandsData = append(b.commandsData, fmt.Sprintf("DROP INDEX IF EXISTS %s", name))
}

func (b *Blueprint) DropForeign(columns ...string) {
	constraintName := fmt.Sprintf("%s_%s_foreign", b.tableName, strings.Join(columns, "_"))
	b.commandsData = append(b.commandsData, fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", b.tableName, constraintName))
}

func (b *Blueprint) DropPrimary() {
	b.commandsData = append(b.commandsData, fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY", b.tableName))
}

func (b *Blueprint) ToSQL() string {
	if len(b.commandsData) > 0 {
		return strings.Join(b.commandsData, ";\n")
	}

	var cols []string
	var constraints []string

	for _, col := range b.columnsData {
		colDef := fmt.Sprintf("%s %s", col.NameField, col.TypeField)

		if !col.NullableField {
			colDef += " " + b.dialectData.WrapNullable(col.NameField)
		}

		if col.PrimaryField {
			colDef += " " + b.dialectData.WrapPrimaryKey(col.NameField)
		}

		if col.UniqueField {
			colDef += " " + b.dialectData.WrapUnique(col.NameField)
		}

		if col.DefaultField != nil {
			colDef += " " + b.dialectData.WrapDefault(col.NameField, *col.DefaultField)
		}

		if col.ReferencesField != nil {
			constraint := b.dialectData.GetForeignKeyConstraint(
				col.NameField,
				col.ReferencesField.TableField,
				col.ReferencesField.ColumnField,
				col.ReferencesField.OnDeleteField,
				col.ReferencesField.OnUpdateField,
			)
			constraints = append(constraints, constraint)
		}

		cols = append(cols, colDef)
	}

	for _, idx := range b.indexesData {
		constraint := b.dialectData.CreateIndex(
			b.tableName,
			idx.NameField,
			idx.ColumnsField,
			idx.UniqueField,
		)
		constraints = append(constraints, constraint)
	}

	allDefs := append(cols, constraints...)
	return fmt.Sprintf(
		"CREATE TABLE %s (\n  %s\n);",
		b.tableName,
		strings.Join(allDefs, ",\n  "),
	)
}
