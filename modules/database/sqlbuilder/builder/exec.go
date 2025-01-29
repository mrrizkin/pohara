package builder

import (
	"database/sql"
	"fmt"
	"strings"
)

func (b *SQLBuilder) Build() (string, []interface{}) {
	if b.rawSQL != "" {
		return b.dialect.ReplacePlaceholders(b.rawSQL), b.rawArgs
	}

	var sql strings.Builder
	var args []interface{}

	switch b.query.operation {
	case "SELECT":
		b.buildSelect(&sql, &args)
	case "INSERT":
		b.buildInsert(&sql, &args)
	case "UPDATE":
		b.buildUpdate(&sql, &args)
	case "DELETE":
		b.buildDelete(&sql, &args)
	}

	return b.dialect.ReplacePlaceholders(sql.String()), args
}

func (b *SQLBuilder) buildSelect(sql *strings.Builder, args *[]interface{}) {
	sql.WriteString("SELECT ")
	if len(b.query.columns) > 0 {
		for i, col := range b.query.columns {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(b.dialect.QuoteIdentifier(col))
		}
	} else {
		sql.WriteString("*")
	}

	sql.WriteString(" FROM ")
	sql.WriteString(b.query.table)

	// Joins
	for _, join := range b.query.joins {
		fmt.Fprintf(sql, " %s JOIN %s ON %s",
			join.joinType, join.table, join.condition)
	}

	// WHERE clauses
	b.buildWhereClauses(sql, args)

	// GROUP BY
	if len(b.query.groupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		for i, col := range b.query.groupBys {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(b.dialect.QuoteIdentifier(col))
		}
	}

	// HAVING
	if b.query.having != "" {
		sql.WriteString(" HAVING ")
		sql.WriteString(b.query.having)
	}

	// ORDER BY
	if len(b.query.orderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(b.query.orderBys, ", "))
	}

	// LIMIT/OFFSET
	if limitOffset := b.dialect.LimitOffset(b.query.limit, b.query.offset); limitOffset != "" {
		sql.WriteString(" ")
		sql.WriteString(limitOffset)
	}

	// RETURNING
	b.buildReturning(sql)
}

func (b *SQLBuilder) buildInsert(sql *strings.Builder, args *[]interface{}) {
	sql.WriteString("INSERT INTO ")
	sql.WriteString(b.query.table)
	sql.WriteString(" (")

	// Columns
	for i, col := range b.query.columns {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(b.dialect.QuoteIdentifier(col))
	}

	sql.WriteString(") VALUES ")

	// Calculate number of values per row
	valuesPerRow := len(b.query.columns)
	if valuesPerRow == 0 {
		return
	}

	// Split values into groups
	for i := 0; i < len(b.query.values); i += valuesPerRow {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString("(")
		for j := 0; j < valuesPerRow; j++ {
			if j > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString("?")
		}
		sql.WriteString(")")
		*args = append(*args, b.query.values[i:i+valuesPerRow]...)
	}

	b.buildReturning(sql)
}

func (b *SQLBuilder) buildUpdate(sql *strings.Builder, args *[]interface{}) {
	sql.WriteString("UPDATE ")
	sql.WriteString(b.query.table)
	sql.WriteString(" SET ")

	// SET clauses
	for i, set := range b.query.sets {
		if i > 0 {
			sql.WriteString(", ")
		}
		fmt.Fprintf(sql, "%s = ?", b.dialect.QuoteIdentifier(set.column))
		*args = append(*args, set.value)
	}

	// WHERE clauses
	b.buildWhereClauses(sql, args)

	b.buildReturning(sql)
}

func (b *SQLBuilder) buildDelete(sql *strings.Builder, args *[]interface{}) {
	sql.WriteString("DELETE FROM ")
	sql.WriteString(b.query.table)
	b.buildWhereClauses(sql, args)
}

func (b *SQLBuilder) buildWhereClauses(sql *strings.Builder, args *[]interface{}) {
	if len(b.query.wheres) > 0 {
		sql.WriteString(" WHERE ")
		for i, where := range b.query.wheres {
			if i > 0 {
				sql.WriteString(" ")
				sql.WriteString(where.operator)
				sql.WriteString(" ")
			}
			sql.WriteString(where.condition)
			*args = append(*args, where.args...)
		}
	}
}

func (b *SQLBuilder) buildReturning(sql *strings.Builder) {
	if len(b.query.returning) > 0 && b.dialect.HasReturning() {
		sql.WriteString(" RETURNING ")
		for i, col := range b.query.returning {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(b.dialect.QuoteIdentifier(col))
		}
	}
}

func (b *SQLBuilder) Exec() (sql.Result, error) {
	sqlStr, args := b.Build()
	if b.tx != nil {
		return b.tx.Exec(sqlStr, args...)
	}
	return b.db.Exec(sqlStr, args...)
}

func (b *SQLBuilder) Get(dest interface{}) error {
	sqlStr, args := b.Build()
	if b.tx != nil {
		return b.tx.Get(dest, sqlStr, args...)
	}
	return b.db.Get(dest, sqlStr, args...)
}

func (b *SQLBuilder) All(dest interface{}) error {
	sqlStr, args := b.Build()
	if b.tx != nil {
		return b.tx.Select(dest, sqlStr, args...)
	}
	return b.db.Select(dest, sqlStr, args...)
}

func (b *SQLBuilder) Count() (int64, error) {
	var count struct {
		Total int64 `db:"total"`
	}

	statement, args := b.Build()
	if b.tx != nil {
		err := b.tx.Get(
			&count,
			fmt.Sprintf("SELECT COUNT(*) as total (%s) AS counting", statement),
			args,
		)
		if err != nil {
			return 0, nil
		}
	} else {
		err := b.db.Get(
			&count,
			fmt.Sprintf("SELECT COUNT(*) as total (%s) AS counting", statement),
			args,
		)
		if err != nil {
			return 0, nil
		}
	}

	return count.Total, nil
}
