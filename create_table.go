package sqls

import (
	"strings"
)

type createTableStatement struct {
	table       []string
	columns     []string
	ifNotExists bool
}

type createTableBuilder struct {
	builder   *sqlBuilder
	statement *createTableStatement
}

func newCreateTableBuilder() *createTableBuilder {
	builder := &createTableBuilder{}
	builder.builder = newSqlBuilder()
	builder.statement = &createTableStatement{}
	return builder
}

func CREATE_TABLE(v string) *createTableBuilder {
	s := newCreateTableBuilder()
	s.statement.table = append(s.statement.table, v)
	return s
}

func (s *createTableBuilder) IF_NOT_EXISTS() *createTableBuilder {
	s.statement.ifNotExists = true
	return s
}

func (s *createTableBuilder) COLUMN(v ...string) *createTableBuilder {
	s.statement.columns = append(s.statement.columns, v...)
	return s
}

func (s *createTableBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *createTableBuilder) String() string {
	var sqlString string
	keyword := "CREATE TABLE"
	if s.statement.ifNotExists {
		keyword += " IF NOT EXISTS"
	}

	sqlString += s.builder.join(keyword, "", s.statement.table, "", "")
	sqlString += s.builder.join("", "(", s.statement.columns, ", ", ")")
	return strings.Trim(sqlString, "\n")
}

func (s *createTableBuilder) Params() []any {
	result := []any{}
	sqlString := s.String()
	matches := s.builder.paramRegexp.FindAllString(sqlString, -1)

	for _, match := range matches {
		result = append(result, s.builder.params[match])
	}

	return result
}
