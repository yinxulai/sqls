package sqls

import (
	"strings"
)

type deleteStatement struct {
	table  []string
	where  []string
	offset string
	limit  string
}

type deleteBuilder struct {
	builder   *sqlBuilder
	statement *deleteStatement
}

func newDeleteBuilder() *deleteBuilder {
	builder := &deleteBuilder{}
	builder.builder = newSqlBuilder()
	builder.statement = &deleteStatement{}
	return builder
}

// 启动删除语句并指定要从中删除的表。
// 一般来说，后面应该跟一个 WHERE 语句！
func DELETE_FROM(v string) *deleteBuilder {
	s := newDeleteBuilder()
	s.statement.table = append(s.statement.table, v)
	return s
}

// 附加一个新的WHERE子句条件，由 AND 串联。
// 可以多次调用，这会导致它每次都将新条件与 AND 串联起来
func (s *deleteBuilder) WHERE(v ...string) *deleteBuilder {
	s.statement.where = append(s.statement.where, v...)
	return s
}

// 用 OR 连接之前和接下来的 WHERE 子句条件。
func (s *deleteBuilder) OR() *deleteBuilder {
	s.WHERE(orString)
	return s
}

// 用 AND 连接之前和接下来的 WHERE 子句条件。
func (s *deleteBuilder) AND() *deleteBuilder {
	s.WHERE(andString)
	return s
}

// 附加一个LIMIT子句。该方法与 SELECT()、UPDATE() 和 DELETE() 一起使用时有效。
// 该方法设计为在使用 SELECT() 时与 OFFSET() 一起使用。
func (s *deleteBuilder) LIMIT(v string) *deleteBuilder {
	s.statement.limit = v
	return s
}

// 附加一个OFFSET子句。该方法与 SELECT() 一起使用时有效。
// 该方法设计为与 LIMIT() 一起使用。
func (s *deleteBuilder) OFFSET(v string) *deleteBuilder {
	s.statement.limit = v
	return s
}

func (s *deleteBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *deleteBuilder) String() string {
	var sqlString string
	sqlString += s.builder.join("DELETE FROM", "", s.statement.table, "", "")
	sqlString += s.builder.join("WHERE", "(", s.statement.where, " AND ", ")")

	if s.statement.offset != "" {
		sqlString += s.builder.join("OFFSET", "", []string{s.statement.offset}, "", "")
	}
	if s.statement.limit != "" {
		sqlString += s.builder.join("LIMIT", "", []string{s.statement.limit}, "", "")
	}
	return strings.Trim(sqlString, "\n")
}

func (s *deleteBuilder) Params() []any {
	result := []any{}
	sqlString := s.String()
	matches := s.builder.paramRegexp.FindAllString(sqlString, -1)

	for _, match := range matches {
		result = append(result, s.builder.params[match])
	}

	return result
}
