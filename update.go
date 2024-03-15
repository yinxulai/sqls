package sqls

import (
	"fmt"
	"strings"
)

type updateStatement struct {
	set            []string
	table          []string
	join           []string
	innerJoin      []string
	outerJoin      []string
	leftOuterJoin  []string
	rightOuterJoin []string
	where          []string
	offset         string
	limit          string
}

type updateBuilder struct {
	builder   *sqlBuilder
	statement *updateStatement
}

func newUpdateBuilder() *updateBuilder {
	builder := &updateBuilder{}
	builder.builder = newSqlBuilder()
	builder.statement = &updateStatement{}
	return builder
}

// 启动更新语句并指定要更新的表。
// 这之后应该是一个或多个 SET() 调用或者 WHERE() 调用。
func UPDATE(table string) *updateBuilder {
	s := newUpdateBuilder()
	s.statement.table = append(s.statement.table, table)
	return s
}

func (s *updateBuilder) SET(key string, value string) *updateBuilder {
	s.statement.set = append(s.statement.set, fmt.Sprintf("%s=%s", key, value))
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *updateBuilder) JOIN(v string) *updateBuilder {
	s.statement.join = append(s.statement.join, v)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *updateBuilder) INNER_JOIN(v string) *updateBuilder {
	s.statement.innerJoin = append(s.statement.innerJoin, v)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *updateBuilder) LEFT_OUTER_JOIN(v string) *updateBuilder {
	s.statement.leftOuterJoin = append(s.statement.leftOuterJoin, v)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *updateBuilder) RIGHT_OUTER_JOIN(v string) *updateBuilder {
	s.statement.rightOuterJoin = append(s.statement.rightOuterJoin, v)
	return s
}

// 附加一个新的WHERE子句条件，由 AND 串联。
// 可以多次调用，这会导致它每次都将新条件与 AND 串联起来
func (s *updateBuilder) WHERE(v string) *updateBuilder {
	s.statement.where = append(s.statement.where, v)
	return s
}

// 用 OR 连接之前和接下来的 WHERE 子句条件。
func (s *updateBuilder) OR() *updateBuilder {
	s.WHERE(orString)
	return s
}

// 用 AND 连接之前和接下来的 WHERE 子句条件。
func (s *updateBuilder) AND() *updateBuilder {
	s.WHERE(andString)
	return s
}

func (s *updateBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *updateBuilder) String() string {
	var sqlString string

	sqlString += s.builder.join("UPDATE", "", s.statement.table, "", "")

	sqlString += s.builder.join("JOIN", "", s.statement.join, " JOIN ", "")
	sqlString += s.builder.join("INNER JOIN", "", s.statement.innerJoin, " INNER JOIN ", "")
	sqlString += s.builder.join("OUTER JOIN", "", s.statement.outerJoin, " OUTER JOIN ", "")
	sqlString += s.builder.join("LEFT OUTER JOIN", "", s.statement.leftOuterJoin, " LEFT OUTER JOIN ", "")
	sqlString += s.builder.join("RIGHT OUTER JOIN", "", s.statement.rightOuterJoin, " RIGHT OUTER JOIN ", "")

	sqlString += s.builder.join("SET", "", s.statement.set, ", ", "")
	sqlString += s.builder.join("WHERE", "(", s.statement.where, " AND ", ")")

	if s.statement.offset != "" {
		sqlString += s.builder.join("OFFSET", "", []string{s.statement.offset}, "", "")
	}
	if s.statement.limit != "" {
		sqlString += s.builder.join("LIMIT", "", []string{s.statement.limit}, "", "")
	}
	return strings.Trim(sqlString, " ")
}

func (s *updateBuilder) Params() []any {
	return s.builder.Params(s.String())
}
