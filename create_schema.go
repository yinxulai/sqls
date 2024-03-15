package sqls

import (
	"strings"
)

type createSchemaStatement struct {
	schema      []string
	ifNotExists bool
}

type createSchemaBuilder struct {
	builder   *sqlBuilder
	statement *createSchemaStatement
}

func newCreateSchemaBuilder() *createSchemaBuilder {
	builder := &createSchemaBuilder{}
	builder.builder = newSqlBuilder()
	builder.statement = &createSchemaStatement{}
	return builder
}

func CREATE_SCHEMA(v string) *createSchemaBuilder {
	s := newCreateSchemaBuilder()
	s.statement.schema = append(s.statement.schema, v)
	return s
}

func (s *createSchemaBuilder) IF_NOT_EXISTS() *createSchemaBuilder {
	s.statement.ifNotExists = true
	return s
}

func (s *createSchemaBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *createSchemaBuilder) String() string {
	var sqlString string
	keyword := "CREATE SCHEMA"
	if s.statement.ifNotExists {
		keyword += " IF NOT EXISTS"
	}

	sqlString += s.builder.join(keyword, "", s.statement.schema, "", "")
	return strings.Trim(sqlString, "\n")
}

func (s *createSchemaBuilder) Params() []any {
	result := []any{}
	sqlString := s.String()
	matches := s.builder.paramRegexp.FindAllString(sqlString, -1)

	for _, match := range matches {
		result = append(result, s.builder.params[match])
	}

	return result
}
