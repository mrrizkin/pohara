package builder

import (
	"github.com/jmoiron/sqlx"
	"github.com/mrrizkin/pohara/modules/database/sqlbuilder/dialect"
)

type SQLBuilder struct {
	dialect dialect.Dialect
	db      *sqlx.DB
	tx      *sqlx.Tx
	query   queryParts
	rawSQL  string
	rawArgs []interface{}
}

type queryParts struct {
	table     string
	operation string
	columns   []string
	values    []interface{}
	sets      []setClause
	wheres    []whereClause
	joins     []joinClause
	orderBys  []string
	groupBys  []string
	having    string
	limit     int
	offset    int
	returning []string
}

type setClause struct {
	column string
	value  interface{}
}

type whereClause struct {
	condition string
	args      []interface{}
	operator  string
}

type joinClause struct {
	joinType  string
	table     string
	condition string
}

func New(dialect dialect.Dialect, db *sqlx.DB) *SQLBuilder {
	return &SQLBuilder{
		dialect: dialect,
		db:      db,
	}
}

func (b *SQLBuilder) WithTx(tx *sqlx.Tx) *SQLBuilder {
	b.tx = tx
	return b
}

func (b *SQLBuilder) Table(name string) *SQLBuilder {
	b.query.table = b.dialect.QuoteIdentifier(name)
	return b
}

func (b *SQLBuilder) Raw(sql string, args ...interface{}) *SQLBuilder {
	b.rawSQL = sql
	b.rawArgs = args
	return b
}

func (b *SQLBuilder) Clone() *SQLBuilder {
	// Deep copy slices
	copySlice := func(src []string) []string {
		dst := make([]string, len(src))
		copy(dst, src)
		return dst
	}

	copyWhereClauses := func(src []whereClause) []whereClause {
		dst := make([]whereClause, len(src))
		for i, w := range src {
			dst[i] = whereClause{
				condition: w.condition,
				operator:  w.operator,
				args:      append([]interface{}{}, w.args...),
			}
		}
		return dst
	}

	return &SQLBuilder{
		dialect: b.dialect,
		db:      b.db,
		tx:      b.tx,
		query: queryParts{
			table:     b.query.table,
			operation: b.query.operation,
			columns:   copySlice(b.query.columns),
			values:    append([]interface{}{}, b.query.values...),
			sets:      append([]setClause{}, b.query.sets...),
			wheres:    copyWhereClauses(b.query.wheres),
			joins:     append([]joinClause{}, b.query.joins...),
			orderBys:  copySlice(b.query.orderBys),
			groupBys:  copySlice(b.query.groupBys),
			having:    b.query.having,
			limit:     b.query.limit,
			offset:    b.query.offset,
			returning: copySlice(b.query.returning),
		},
		rawSQL:  b.rawSQL,
		rawArgs: append([]interface{}{}, b.rawArgs...),
	}
}
