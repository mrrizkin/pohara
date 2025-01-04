package ports

import "github.com/mrrizkin/pohara/internal/common/sql"

type Pagination struct {
	Limit  sql.Int64Nullable
	Offset sql.Int64Nullable
	Sort   sql.StringNullable
}

type FindResult struct {
	Data      interface{}
	Total     int64
	TotalPage sql.Int64Nullable
	Page      sql.Int64Nullable
	Limit     sql.Int64Nullable
}

type Database interface {
	GetDB() interface{}
	Create(value interface{}) error
	First(dest interface{}, conds ...interface{}) error
	FindAll(dest interface{}, conds ...interface{}) error
	Find(dest interface{}, p Pagination, conds ...interface{}) (*FindResult, error)
	Count(model interface{}, conds ...interface{}) (int64, error)
	Save(value interface{}) error
	Delete(value interface{}, conds ...interface{}) error
	Transaction(fn func(tx Database) error) error
	Raw(dest interface{}, statement string, values ...interface{}) error
	Exec(statement string, values ...interface{}) error
}
