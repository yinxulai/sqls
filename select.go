package sqls

import (
	"strings"
)

type selectStatement struct {
	selects        []string
	distinct       bool
	table          []string
	join           []string
	innerJoin      []string
	outerJoin      []string
	leftOuterJoin  []string
	rightOuterJoin []string
	where          []string
	having         []string
	groupBy        []string
	orderBy        []string
	offset         string
	limit          string

	// 标记当前是处于 where 还是 having
	// 以便 or、and 正确的 push 到对应的 list 去
	lastContext string
}

type selectBuilder struct {
	builder   *sqlBuilder
	statement *selectStatement
}

func newSelectBuilder() *selectBuilder {
	builder := &selectBuilder{}
	builder.builder = newSqlBuilder()
	builder.statement = &selectStatement{}
	return builder
}

// 开始或附加到 SELECT 子句。
// 可以多次调用，并且参数将附加到SELECT子句中。
// 这些参数通常是逗号分隔的列和别名列表，但可以是驱动程序可接受的任何内容。
func SELECT(v ...string) *selectBuilder {
	s := newSelectBuilder()
	s.statement.selects = append(s.statement.selects, v...)
	return s
}

// 开始或附加到SELECT子句，还将DISTINCT关键字添加到生成的查询中。
// 可以多次调用，并且参数将附加到SELECT子句中。
// 这些参数通常是逗号分隔的列和别名列表，但可以是驱动程序可接受的任何内容。
func SELECT_DISTINCT(v ...string) *selectBuilder {
	s := newSelectBuilder()
	s.statement.distinct = true
	s.statement.selects = append(s.statement.selects, v...)
	return s
}

func (s *selectBuilder) SELECT(v ...string) *selectBuilder {
	s.statement.selects = append(s.statement.selects, v...)
	return s
}

func (s *selectBuilder) SELECT_DISTINCT(v ...string) *selectBuilder {
	s.statement.distinct = true
	s.statement.selects = append(s.statement.selects, v...)
	return s
}

// 开始或附加到FROM子句。可以多次调用，并且参数将附加到FROM子句中。
// 参数通常是表名和别名，或者驱动程序可接受的任何内容。
func (s *selectBuilder) FROM(v string) *selectBuilder {
	s.statement.table = append(s.statement.table, v)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *selectBuilder) JOIN(v string) *selectBuilder {
	s.statement.join = append(s.statement.join, v)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *selectBuilder) INNER_JOIN(v string) *selectBuilder {
	s.statement.innerJoin = append(s.statement.innerJoin, v)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *selectBuilder) LEFT_OUTER_JOIN(v string) *selectBuilder {
	s.statement.leftOuterJoin = append(s.statement.leftOuterJoin, v)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *selectBuilder) RIGHT_OUTER_JOIN(v string) *selectBuilder {
	s.statement.rightOuterJoin = append(s.statement.rightOuterJoin, v)
	return s
}

// 附加一个新的WHERE子句条件，由 AND 串联。
// 可以多次调用，这会导致它每次都将新条件与 AND 串联起来
func (s *selectBuilder) WHERE(v string) *selectBuilder {
	s.statement.where = append(s.statement.where, v)
	s.statement.lastContext = "WHERE"
	return s
}

// 用 OR 连接之前和接下来的 WHERE 子句条件。
func (s *selectBuilder) OR() *selectBuilder {
	if s.statement.lastContext == "WHERE" {
		s.WHERE(orString)
	}
	if s.statement.lastContext == "HAVING" {
		s.HAVING(orString)
	}
	return s
}

// 用 AND 连接之前和接下来的 WHERE 子句条件。
func (s *selectBuilder) AND() *selectBuilder {
	if s.statement.lastContext == "WHERE" {
		s.WHERE(andString)
	}
	if s.statement.lastContext == "HAVING" {
		s.HAVING(andString)
	}
	return s
}

// 追加一个新的 GROUP BY 子句元素，并用逗号连接。
// 可以多次调用，这会导致它每次都用逗号连接新条件。
func (s *selectBuilder) GROUP_BY(v string) *selectBuilder {
	s.statement.groupBy = append(s.statement.groupBy, v)
	return s
}

// 追加一个新的 HAVING 子句条件，并通过 AND 连接。
// 可以多次调用，这会导致它每次都将新条件与 串联起来AND。使用 OR() 来分割 OR。
func (s *selectBuilder) HAVING(v string) *selectBuilder {
	s.statement.having = append(s.statement.having, v)
	s.statement.lastContext = "HAVING"
	return s
}

// 追加一个新的ORDER BY子句元素，并用逗号连接。可以多次调用，这会导致它每次都用逗号连接新条件。
func (s *selectBuilder) ORDER_BY(v string) *selectBuilder {
	s.statement.orderBy = append(s.statement.orderBy, v)
	return s
}

// 附加一个LIMIT子句。该方法与 SELECT()、UPDATE() 和 DELETE() 一起使用时有效。
// 该方法设计为在使用 SELECT() 时与 OFFSET() 一起使用。
func (s *selectBuilder) LIMIT(v string) *selectBuilder {
	s.statement.limit = v
	return s
}

// 附加一个OFFSET子句。该方法与 SELECT() 一起使用时有效。
// 该方法设计为与 LIMIT() 一起使用。
func (s *selectBuilder) OFFSET(v string) *selectBuilder {
	s.statement.offset = v
	return s
}

func (s *selectBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *selectBuilder) String() string {
	var sqlString string

	if s.statement.distinct {
		sqlString += s.builder.join("SELECT DISTINCT", "", s.statement.selects, ", ", "")
	} else {
		sqlString += s.builder.join("SELECT", "", s.statement.selects, ", ", "")
	}

	sqlString += s.builder.join("FROM", "", s.statement.table, ", ", "")

	sqlString += s.builder.join("JOIN", "", s.statement.join, "\nJOIN ", "")
	sqlString += s.builder.join("INNER JOIN", "", s.statement.innerJoin, "\nINNER JOIN ", "")
	sqlString += s.builder.join("OUTER JOIN", "", s.statement.outerJoin, "\nOUTER JOIN ", "")
	sqlString += s.builder.join("LEFT OUTER JOIN", "", s.statement.leftOuterJoin, "\nLEFT OUTER JOIN ", "")
	sqlString += s.builder.join("RIGHT OUTER JOIN", "", s.statement.rightOuterJoin, "\nRIGHT OUTER JOIN ", "")

	sqlString += s.builder.join("WHERE", "(", s.statement.where, " AND ", ")")
	sqlString += s.builder.join("GROUP BY", "", s.statement.groupBy, ", ", "")
	sqlString += s.builder.join("HAVING", "(", s.statement.having, " AND ", ")")
	sqlString += s.builder.join("ORDER BY", "", s.statement.orderBy, ", ", "")
	if s.statement.offset != "" {
		sqlString += s.builder.join("OFFSET", "", []string{s.statement.offset}, "", "")
	}
	if s.statement.limit != "" {
		sqlString += s.builder.join("LIMIT", "", []string{s.statement.limit}, "", "")
	}
	return strings.Trim(sqlString, "\n")
}

func (s *selectBuilder) Params() []any {
	return s.builder.Params(s.String())
}
