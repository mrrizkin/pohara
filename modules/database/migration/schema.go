package migration

import "strings"

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

	if len(bp.CompositePrimary) > 0 {
		s.statements = append(
			s.statements,
			s.dialect.CreateCompositePrimarySQL(bp.CompositePrimary...),
		)
	}
}

func (s *Schema) CreateNotExist(table string, callback func(*Blueprint)) {
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

		if col.Index != "" {
			bp.Indexes = append(bp.Indexes, index{
				Name:    col.Index,
				Columns: []string{col.Name},
				Unique:  col.Unique,
			})
		}
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

func (s *Schema) DropExist(table string) {
	s.statements = append(s.statements, s.dialect.DropTableIfExistSQL(table))
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

	if col.Foreign.Table != "" && col.Foreign.Column != "" {
		parts = append(parts, "REFERENCES "+col.Foreign.Table+"("+col.Foreign.Column+")")
		if col.Foreign.OnDelete != "" {
			parts = append(parts, "ON DELETE "+col.Foreign.OnDelete)
		}
	}

	return strings.Join(parts, " ")
}

func (s *Schema) Raw(statement string) {
	s.statements = append(s.statements, statement)
}

func (s *Schema) statement() []string {
	return s.statements
}
