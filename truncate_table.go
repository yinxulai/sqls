package sqls

import (
	"strings"
)

type truncateTableStatement struct {
	table string
}

type truncateTableBuilder struct {
	builder   *sqlBuilder
	statement *truncateTableStatement
}

func newTruncateTableBuilder() *truncateTableBuilder {
	builder := &truncateTableBuilder{}
	builder.builder = newSqlBuilder()
	builder.statement = &truncateTableStatement{}
	return builder
}

func TRUNCATE_TABLE(v string) *truncateTableBuilder {
	s := newTruncateTableBuilder()
	s.statement.table = v
	return s
}

func (s *truncateTableBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *truncateTableBuilder) String() string {
	var sqlString string
	sqlString += s.builder.join("TRUNCATE TABLE", "", []string{s.statement.table}, "", "")
	return strings.Trim(sqlString, "\n")
}

func (s *truncateTableBuilder) Params() []any {
	return s.builder.Params(s.String())
}
