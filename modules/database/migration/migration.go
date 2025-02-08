package migration

import (
	"strings"
)

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
	RenameTableSQL(old, new string) string
	CreateIndexSQL(table, name string, columns []string, unique bool) string
	DropIndexSQL(table, name string) string
}

// Column definition
type Column struct {
	Name      string
	Type      string
	Nullable  bool
	Default   string
	Unique    bool
	Modify    bool
	After     string
	Length    int
	Precision int
	Scale     int
	IsPrimary bool
}

// Blueprint for table definition
type Blueprint struct {
	TableName  string
	Columns    []Column
	Indexes    []index
	Renames    []rename
	Drops      []string
	Modifies   []Column
	IsCreating bool
	dialect    Dialect
}

type index struct {
	Columns []string
	Unique  bool
	Name    string
}

type rename struct {
	From string
	To   string
}

type Schema struct {
	dialect    Dialect
	statements []string
}

func NewSchema(dialect Dialect) *Schema {
	return &Schema{dialect: dialect}
}

// Create a new table
func (s *Schema) Create(table string, callback func(*Blueprint)) {
	bp := &Blueprint{TableName: table, IsCreating: true, dialect: s.dialect}
	callback(bp)

	columns := []string{}
	for _, col := range bp.Columns {
		columns = append(columns, s.buildColumn(col))
	}

	sql := s.dialect.CreateTableSQL(table, columns)
	s.statements = append(s.statements, sql)

	for _, idx := range bp.Indexes {
		s.statements = append(
			s.statements,
			s.dialect.CreateIndexSQL(table, idx.Name, idx.Columns, idx.Unique),
		)
	}
}

func (s *Schema) CreateIfNotExist(table string, callback func(*Blueprint)) {
	bp := &Blueprint{TableName: table, IsCreating: true, dialect: s.dialect}
	callback(bp)

	columns := []string{}
	for _, col := range bp.Columns {
		columns = append(columns, s.buildColumn(col))
	}

	sql := s.dialect.CreateTableIfNotExistSQL(table, columns)
	s.statements = append(s.statements, sql)

	for _, idx := range bp.Indexes {
		s.statements = append(
			s.statements,
			s.dialect.CreateIndexSQL(table, idx.Name, idx.Columns, idx.Unique),
		)
	}
}

// Modify existing table
func (s *Schema) Table(table string, callback func(*Blueprint)) {
	bp := &Blueprint{TableName: table, dialect: s.dialect}
	callback(bp)

	// Handle column modifications
	for _, modify := range bp.Modifies {
		sql := s.dialect.ModifyColumnSQL(table, s.buildColumn(modify))
		s.statements = append(s.statements, sql)
	}

	// Handle new columns
	for _, col := range bp.Columns {
		sql := s.dialect.AddColumnSQL(table, s.buildColumn(col))
		s.statements = append(s.statements, sql)
	}

	// Handle renames
	for _, rn := range bp.Renames {
		sql := s.dialect.RenameColumnSQL(table, rn.From, rn.To)
		s.statements = append(s.statements, sql)
	}

	// Handle drops
	for _, drop := range bp.Drops {
		sql := s.dialect.DropColumnSQL(table, drop)
		s.statements = append(s.statements, sql)
	}

	// Handle index drops
	for _, idx := range bp.Indexes {
		if idx.Name != "" {
			s.statements = append(s.statements, s.dialect.DropIndexSQL(table, idx.Name))
		}
	}
}

// Drop table
func (s *Schema) Drop(table string) {
	s.statements = append(s.statements, s.dialect.DropTableSQL(table))
}

// Rename table
func (s *Schema) Rename(old, new string) {
	s.statements = append(s.statements, s.dialect.RenameTableSQL(old, new))
}

// Build column definition
func (s *Schema) buildColumn(col Column) string {
	parts := []string{col.Name, col.Type}

	if !col.Nullable {
		parts = append(parts, s.dialect.GetNullable(col.Name, false))
	}

	if col.Default != "" {
		parts = append(parts, s.dialect.GetDefault(col.Name, col.Default))
	}

	if col.Unique && !col.IsPrimary {
		parts = append(parts, s.dialect.GetUnique(col.Name))
	}
	if col.IsPrimary {
		parts = append(parts, s.dialect.GetPrimaryKey())
	}

	if col.After != "" {
		parts = append(parts, "AFTER", col.After)
	}

	return strings.Join(parts, " ")
}

func (s *Schema) Raw(statement string) {
	s.statements = append(s.statements, statement)
}

func (s *Schema) statement() []string {
	return s.statements
}

// Blueprint methods
func (b *Blueprint) ID() {
	col := Column{
		Name:      "id",
		Type:      b.dialect.GetBigIntegerType(),
		IsPrimary: true,
	}
	b.Columns = append(b.Columns, col)
}

func (b *Blueprint) String(name string, length int) *ColumnBuilder {
	col := Column{
		Name:   name,
		Type:   b.dialect.GetStringType(length),
		Length: length,
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) Text(name string) *ColumnBuilder {
	col := Column{
		Name: name,
		Type: b.dialect.GetTextType(),
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) Json(name string) *ColumnBuilder {
	col := Column{
		Name: name,
		Type: b.dialect.GetJsonType(),
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) Integer(name string) *ColumnBuilder {
	col := Column{
		Name: name,
		Type: b.dialect.GetIntegerType(),
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) BigInteger(name string) *ColumnBuilder {
	col := Column{
		Name: name,
		Type: b.dialect.GetBigIntegerType(),
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) Boolean(name string) *ColumnBuilder {
	col := Column{
		Name: name,
		Type: b.dialect.GetBooleanType(),
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) Enum(name string, values []string) *ColumnBuilder {
	col := Column{
		Name: name,
		Type: b.dialect.GetEnumType(values),
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) Timestamp(name string) *ColumnBuilder {
	col := Column{
		Name: name,
		Type: b.dialect.GetTimestampType(),
	}
	b.Columns = append(b.Columns, col)
	return &ColumnBuilder{col: &b.Columns[len(b.Columns)-1]}
}

func (b *Blueprint) Timestamps() {
	b.Columns = append(b.Columns,
		Column{
			Name:     "created_at",
			Type:     b.dialect.GetTimestampType(),
			Nullable: true,
		},
		Column{
			Name:     "updated_at",
			Type:     b.dialect.GetTimestampType(),
			Nullable: true,
		},
	)
}

func (b *Blueprint) RenameColumn(old, new string) {
	b.Renames = append(b.Renames, rename{From: old, To: new})
}

func (b *Blueprint) DropColumn(name string) {
	b.Drops = append(b.Drops, name)
}

func (b *Blueprint) ModifyColumn(column Column) {
	column.Modify = true
	b.Modifies = append(b.Modifies, column)
}

// Column builder for fluent interface
type ColumnBuilder struct {
	col *Column
}

func (cb *ColumnBuilder) Primary() *ColumnBuilder {
	cb.col.IsPrimary = true
	return cb
}

func (cb *ColumnBuilder) Nullable() *ColumnBuilder {
	cb.col.Nullable = true
	return cb
}

func (cb *ColumnBuilder) Default(value string) *ColumnBuilder {
	cb.col.Default = value
	return cb
}

func (cb *ColumnBuilder) Unique() *ColumnBuilder {
	cb.col.Unique = true
	return cb
}

func (cb *ColumnBuilder) After(column string) *ColumnBuilder {
	cb.col.After = column
	return cb
}
