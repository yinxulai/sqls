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
	return strings.Trim(sqlString, " ")
}

func (s *createSchemaBuilder) Params() []any {
	return s.builder.Params(s.String())
}
