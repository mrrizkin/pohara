package sql

type whereBuilder struct {
	count int
	where string
	args  []interface{}
}

func Where() *whereBuilder {
	return &whereBuilder{
		count: 0,
		where: "",
		args:  make([]interface{}, 0),
	}
}

func (wb *whereBuilder) And(where string, args ...interface{}) {
	if wb.count != 0 {
		wb.where += " AND"
	}

	wb.where += " " + where
	wb.args = append(wb.args, args...)
	wb.count++
}

func (wb *whereBuilder) Or(where string, args ...interface{}) {
	if wb.count != 0 {
		wb.where += " OR"
	}

	wb.where += " " + where
	wb.args = append(wb.args, args...)
	wb.count++
}

func (wb *whereBuilder) Get() (string, []interface{}) {
	return wb.where, wb.args
}

func (wb *whereBuilder) Cond() []interface{} {
	conds := []interface{}{wb.where}
	return append(conds, wb.args...)
}
