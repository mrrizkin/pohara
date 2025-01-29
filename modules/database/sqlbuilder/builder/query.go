package builder

import (
	"fmt"
	"reflect"
)

func (b *SQLBuilder) Select(columns ...string) *SQLBuilder {
	b.query.operation = "SELECT"
	b.query.columns = columns
	return b
}

func (b *SQLBuilder) Insert(data ...interface{}) *SQLBuilder {
	b.query.operation = "INSERT"
	if len(data) > 0 {
		b.Values(data[0])
	}
	return b
}

func (b *SQLBuilder) Update() *SQLBuilder {
	b.query.operation = "UPDATE"
	return b
}

func (b *SQLBuilder) Delete() *SQLBuilder {
	b.query.operation = "DELETE"
	return b
}

func (b *SQLBuilder) Values(data interface{}) *SQLBuilder {
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		typ := rv.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tag := field.Tag.Get("db")
			if tag == "" {
				continue
			}
			b.query.columns = append(b.query.columns, tag)
			b.query.values = append(b.query.values, rv.Field(i).Interface())
		}
	case reflect.Map:
		for k, v := range data.(map[string]interface{}) {
			b.query.columns = append(b.query.columns, k)
			b.query.values = append(b.query.values, v)
		}
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			if i == 0 {
				b.Values(s.Index(i).Interface())
			} else {
				b.AddValues(s.Index(i).Interface())
			}
		}
	}
	return b
}

func (b *SQLBuilder) AddValues(data interface{}) *SQLBuilder {
	if len(b.query.columns) == 0 {
		panic("AddValues called without initial Values call")
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	var newValues []interface{}

	switch rv.Kind() {
	case reflect.Struct:
		typ := rv.Type()
		columnIndex := make(map[string]int)
		for i, col := range b.query.columns {
			columnIndex[col] = i
		}

		// Verify all columns exist in struct
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tag := field.Tag.Get("db")
			if tag != "" {
				if _, exists := columnIndex[tag]; !exists {
					panic(fmt.Sprintf("column %s not in initial insert columns", tag))
				}
			}
		}

		// Collect values in column order
		for _, col := range b.query.columns {
			found := false
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				if field.Tag.Get("db") == col {
					newValues = append(newValues, rv.Field(i).Interface())
					found = true
					break
				}
			}
			if !found {
				panic(fmt.Sprintf("missing column %s in struct", col))
			}
		}

	case reflect.Map:
		m, ok := data.(map[string]interface{})
		if !ok {
			panic("AddValues requires map[string]interface{}")
		}

		// Verify all columns exist in map
		for _, col := range b.query.columns {
			if _, exists := m[col]; !exists {
				panic(fmt.Sprintf("missing column %s in map", col))
			}
		}

		// Collect values in column order
		for _, col := range b.query.columns {
			newValues = append(newValues, m[col])
		}

	default:
		panic("unsupported type for AddValues")
	}

	b.query.values = append(b.query.values, newValues...)
	return b
}

func (b *SQLBuilder) Set(column string, value interface{}) *SQLBuilder {
	b.query.sets = append(b.query.sets, setClause{
		column: column,
		value:  value,
	})
	return b
}

func (b *SQLBuilder) Where(condition string, args ...interface{}) *SQLBuilder {
	return b.addWhereClause("AND", condition, args...)
}

func (b *SQLBuilder) OrWhere(condition string, args ...interface{}) *SQLBuilder {
	return b.addWhereClause("OR", condition, args...)
}

func (b *SQLBuilder) addWhereClause(operator, condition string, args ...interface{}) *SQLBuilder {
	b.query.wheres = append(b.query.wheres, whereClause{
		operator:  operator,
		condition: condition,
		args:      args,
	})
	return b
}

func (b *SQLBuilder) Join(joinType, table, condition string) *SQLBuilder {
	b.query.joins = append(b.query.joins, joinClause{
		joinType:  joinType,
		table:     b.dialect.QuoteIdentifier(table),
		condition: condition,
	})
	return b
}

func (b *SQLBuilder) OrderBy(column string, direction string) *SQLBuilder {
	b.query.orderBys = append(b.query.orderBys,
		b.dialect.QuoteIdentifier(column)+" "+direction)
	return b
}

func (b *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	b.query.groupBys = append(b.query.groupBys, columns...)
	return b
}

func (b *SQLBuilder) Having(condition string) *SQLBuilder {
	b.query.having = condition
	return b
}

func (b *SQLBuilder) Limit(limit int) *SQLBuilder {
	b.query.limit = limit
	return b
}

func (b *SQLBuilder) Offset(offset int) *SQLBuilder {
	b.query.offset = offset
	return b
}

func (b *SQLBuilder) Returning(columns ...string) *SQLBuilder {
	b.query.returning = columns
	return b
}
