package migration

// Blueprint for table definition
type Blueprint struct {
	TableName        string
	Columns          []Column
	Indexes          []index
	Renames          []rename
	Drops            []string
	Modifies         []Column
	CompositePrimary []string
	IsCreating       bool
	dialect          Dialect
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

// Blueprint methods
func (b *Blueprint) ID() {
	col := Column{
		Name:      "id",
		Type:      b.dialect.GetBigIntegerType(),
		IsPrimary: true,
	}
	b.Columns = append(b.Columns, col)
}

func (b *Blueprint) Primary(columns ...string) {
	b.CompositePrimary = append(b.CompositePrimary, columns...)
}

func (b *Blueprint) Index(name string, unique bool, columns ...string) {
	b.Indexes = append(b.Indexes, index{
		Name:    name,
		Unique:  unique,
		Columns: columns,
	})
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
