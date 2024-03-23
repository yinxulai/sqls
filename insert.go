package sqls

import (
	"strings"
)

type insertStatement struct {
	table       []string
	columns     []string
	values      []string
	onConflict  []string
	doUpdateSet []string
}

type insertBuilder struct {
	builder   *sqlBuilder
	statement *insertStatement
}

func newInsertBuilder() *insertBuilder {
	builder := &insertBuilder{}
	builder.builder = newSqlBuilder()
	builder.statement = &insertStatement{}
	return builder
}

// 启动插入语句并指定要插入的表。
// 此后应跟一个或多个 VALUES()
func INSERT_INTO(v string) *insertBuilder {
	s := newInsertBuilder()
	s.statement.table = append(s.statement.table, v)
	return s
}

// 附加到插入语句。第一个参数是要插入的列，第二个参数是值。
func (s *insertBuilder) VALUES(key string, value string) *insertBuilder {
	s.statement.columns = append(s.statement.columns, key)
	s.statement.values = append(s.statement.values, value)
	return s
}

func (s *insertBuilder) ON_CONFLICT(v string) *insertBuilder {
	s.statement.onConflict = append(s.statement.onConflict, v)
	return s
}

func (s *insertBuilder) DO_UPDATE_SET(key string, value string) *insertBuilder {
	s.statement.doUpdateSet = append(s.statement.onConflict, key+"="+value)
	return s
}

func (s *insertBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *insertBuilder) String() string {
	var sqlString string
	sqlString += s.builder.join("INSERT INTO", "", s.statement.table, "", "")
	sqlString += s.builder.join("", "(", s.statement.columns, ", ", ")")
	sqlString += s.builder.join("VALUES", "(", s.statement.values, ", ", ")")

	if len(s.statement.onConflict) > 0 && len(s.statement.doUpdateSet) > 0 {
		sqlString += s.builder.join("ON CONFLICT", "(", s.statement.onConflict, ", ", ")")
		sqlString += s.builder.join("DO UPDATE SET", "", s.statement.onConflict, ", ", "")
	}

	return strings.Trim(sqlString, " ")
}

func (s *insertBuilder) Params() []any {
	return s.builder.Params(s.String())
}
