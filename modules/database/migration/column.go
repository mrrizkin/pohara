package migration

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
	Index     string
	IsPrimary bool
	Foreign   Foreign
}

type Foreign struct {
	Table    string
	Column   string
	OnDelete string
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

func (cb *ColumnBuilder) Index(column string) *ColumnBuilder {
	cb.col.Index = column
	return cb
}

func (cb *ColumnBuilder) Foreign(table, column string) *ColumnBuilder {
	cb.col.Foreign.Table = table
	cb.col.Foreign.Column = column
	return cb
}

func (cb *ColumnBuilder) OnDelete(action string) *ColumnBuilder {
	cb.col.Foreign.OnDelete = action
	return cb
}
