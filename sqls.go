package sqls

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const (
	DELETE        = "DELETE"
	SELECT        = "SELECT"
	UPDATE        = "UPDATE"
	INSERT_INTO   = "INSERT INTO"
	CREATE_TABLE  = "CREATE TABLE"
	CREATE_SCHEMA = "CREATE SCHEMA"
)

type statement struct {
	set            []string
	schema         string
	selects        []string
	distinct       bool
	ifNotExists    bool
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
	columns        []string
	values         []string
	offset         string
	limit          string

	// 标记当前是处于 where 还是 having
	// 以便 or、and 正确的 push 到对应的 list 去
	lastContext string
}

type sqlBuilder struct {
	class       string
	statement   *statement
	params      map[string]any
	paramRegexp *regexp.Regexp
}

func Begin() *sqlBuilder {
	builder := &sqlBuilder{}
	builder.statement = &statement{}
	builder.params = make(map[string]any)
	builder.paramRegexp = regexp.MustCompile("$([0-9]+)")
	return builder
}

const orString = ") OR ("
const andString = ") AND ("

// 拼接 sql 的工具方法
func (s *sqlBuilder) join(keyword string, open string, fields []string, sep string, close string) string {
	if len(fields) == 0 {
		return ""
	}

	var body = ""
	for index, field := range fields {
		if (index > 0) && (index != len(field)-1) {
			// 在 or、and 的前后都不应该添加 sep
			if !slices.Contains[[]string]([]string{orString, andString}, field) {
				if !slices.Contains[[]string]([]string{orString, andString}, fields[index-1]) {
					body += sep
				}
			}
		}

		body += field
	}

	if keyword != "" {
		return keyword + " " + open + body + close + "\n"
	}

	return open + body + close + "\n"
}

// 开始或附加到 SELECT 子句。
// 可以多次调用，并且参数将附加到SELECT子句中。
// 这些参数通常是逗号分隔的列和别名列表，但可以是驱动程序可接受的任何内容。
func (s *sqlBuilder) SELECT(v ...string) *sqlBuilder {
	s.class = SELECT
	s.statement.selects = append(s.statement.selects, v...)
	return s
}

// 开始或附加到SELECT子句，还将DISTINCT关键字添加到生成的查询中。
// 可以多次调用，并且参数将附加到SELECT子句中。
// 这些参数通常是逗号分隔的列和别名列表，但可以是驱动程序可接受的任何内容。
func (s *sqlBuilder) SELECT_DISTINCT(v ...string) *sqlBuilder {
	s.class = SELECT
	s.statement.distinct = true
	s.statement.selects = append(s.statement.selects, v...)
	return s
}

// 开始或附加到FROM子句。可以多次调用，并且参数将附加到FROM子句中。
// 参数通常是表名和别名，或者驱动程序可接受的任何内容。
func (s *sqlBuilder) FROM(v ...string) *sqlBuilder {
	s.statement.table = append(s.statement.table, v...)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *sqlBuilder) JOIN(v ...string) *sqlBuilder {
	s.statement.join = append(s.statement.join, v...)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *sqlBuilder) INNER_JOIN(v ...string) *sqlBuilder {
	s.statement.innerJoin = append(s.statement.innerJoin, v...)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *sqlBuilder) LEFT_OUTER_JOIN(v ...string) *sqlBuilder {
	s.statement.leftOuterJoin = append(s.statement.leftOuterJoin, v...)
	return s
}

// JOIN 根据调用的方法添加适当类型的新子句。
// 该参数可以包括由列和连接条件组成的标准连接。
func (s *sqlBuilder) RIGHT_OUTER_JOIN(v ...string) *sqlBuilder {
	s.statement.rightOuterJoin = append(s.statement.rightOuterJoin, v...)
	return s
}

// 附加一个新的WHERE子句条件，由 AND 串联。
// 可以多次调用，这会导致它每次都将新条件与 AND 串联起来
func (s *sqlBuilder) WHERE(v ...string) *sqlBuilder {
	s.statement.where = append(s.statement.where, v...)
	s.statement.lastContext = "WHERE"
	return s
}

// 用 OR 连接之前和接下来的 WHERE 子句条件。
func (s *sqlBuilder) OR() *sqlBuilder {
	if s.statement.lastContext == "WHERE" {
		s.WHERE(orString)
	}
	if s.statement.lastContext == "HAVING" {
		s.HAVING(orString)
	}
	return s
}

// 用 AND 连接之前和接下来的 WHERE 子句条件。
func (s *sqlBuilder) AND() *sqlBuilder {
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
func (s *sqlBuilder) GROUP_BY(v ...string) *sqlBuilder {
	s.statement.groupBy = append(s.statement.groupBy, v...)
	return s
}

// 追加一个新的 HAVING 子句条件，并通过 AND 连接。
// 可以多次调用，这会导致它每次都将新条件与 串联起来AND。使用 OR() 来分割 OR。
func (s *sqlBuilder) HAVING(v ...string) *sqlBuilder {
	s.statement.having = append(s.statement.having, v...)
	s.statement.lastContext = "HAVING"
	return s
}

// 追加一个新的ORDER BY子句元素，并用逗号连接。可以多次调用，这会导致它每次都用逗号连接新条件。
func (s *sqlBuilder) ORDER_BY(v ...string) *sqlBuilder {
	s.statement.orderBy = append(s.statement.orderBy, v...)
	return s
}

// 附加一个LIMIT子句。该方法与 SELECT()、UPDATE() 和 DELETE() 一起使用时有效。
// 该方法设计为在使用 SELECT() 时与 OFFSET() 一起使用。
func (s *sqlBuilder) LIMIT(v string) *sqlBuilder {
	s.statement.limit = v
	return s
}

// 附加一个OFFSET子句。该方法与 SELECT() 一起使用时有效。
// 该方法设计为与 LIMIT() 一起使用。
func (s *sqlBuilder) OFFSET(v string) *sqlBuilder {
	s.statement.limit = v
	return s
}

// 启动删除语句并指定要从中删除的表。
// 一般来说，后面应该跟一个 WHERE 语句！
func (s *sqlBuilder) DELETE_FROM(v string) *sqlBuilder {
	s.class = DELETE
	s.statement.table = append(s.statement.table, v)
	return s
}

// 启动插入语句并指定要插入的表。
// 此后应跟一个或多个 VALUES()
func (s *sqlBuilder) INSERT_INTO(v string) *sqlBuilder {
	s.class = INSERT_INTO
	s.statement.table = append(s.statement.table, v)
	return s
}

// 附加到插入语句。第一个参数是要插入的列，第二个参数是值。
func (s *sqlBuilder) VALUES(key string, value string) *sqlBuilder {
	s.statement.columns = append(s.statement.columns, key)
	s.statement.values = append(s.statement.values, value)
	return s
}

// 启动更新语句并指定要更新的表。
// 这之后应该是一个或多个 SET() 调用或者 WHERE() 调用。
func (s *sqlBuilder) UPDATE(v string) *sqlBuilder {
	s.class = UPDATE
	s.statement.table = append(s.statement.table, v)
	return s
}

// 需要配合 UPDATE 生效
func (s *sqlBuilder) SET(v ...string) *sqlBuilder {
	s.statement.set = append(s.statement.set, v...)
	return s
}

// 创建数据表
// 后面可以接 IF_NOT_EXISTS
func (s *sqlBuilder) CREATE_TABLE(t string, fields ...string) *sqlBuilder {
	s.class = CREATE_TABLE
	s.statement.table = append(s.statement.table, t)
	s.statement.columns = append(s.statement.columns, fields...)
	return s
}

func (s *sqlBuilder) CREATE_SCHEMA(t string) *sqlBuilder {
	s.class = CREATE_SCHEMA
	s.statement.schema = t
	return s
}

func (s *sqlBuilder) IF_NOT_EXISTS() *sqlBuilder {
	s.statement.ifNotExists = true
	return s
}

func (s *sqlBuilder) IF_EXISTS() *sqlBuilder {
	s.statement.ifNotExists = false
	return s
}

func (s *sqlBuilder) Param(v any) string {
	paramsIndex := len(s.params) + 1
	key := fmt.Sprintf("$%d", paramsIndex)
	s.params[key] = v
	return key
}

func (s *sqlBuilder) String() string {
	var sqlString string

	if s.class == SELECT {
		if s.statement.distinct {
			sqlString += s.join("SELECT DISTINCT", "", s.statement.selects, ", ", "")
		} else {
			sqlString += s.join("SELECT", "", s.statement.selects, ", ", "")
		}

		sqlString += s.join("FROM", "", s.statement.table, ", ", "")

		sqlString += s.join("JOIN", "", s.statement.join, "\nJOIN ", "")
		sqlString += s.join("INNER JOIN", "", s.statement.innerJoin, "\nINNER JOIN ", "")
		sqlString += s.join("OUTER JOIN", "", s.statement.outerJoin, "\nOUTER JOIN ", "")
		sqlString += s.join("LEFT OUTER JOIN", "", s.statement.leftOuterJoin, "\nLEFT OUTER JOIN ", "")
		sqlString += s.join("RIGHT OUTER JOIN", "", s.statement.rightOuterJoin, "\nRIGHT OUTER JOIN ", "")

		sqlString += s.join("WHERE", "(", s.statement.where, " AND ", ")")
		sqlString += s.join("GROUP BY", "", s.statement.groupBy, ", ", "")
		sqlString += s.join("HAVING", "(", s.statement.having, " AND ", ")")
		sqlString += s.join("ORDER BY", "", s.statement.orderBy, ", ", "")
		if s.statement.offset != "" {
			sqlString += s.join("OFFSET", "", []string{s.statement.offset}, "", "")
		}
		if s.statement.limit != "" {
			sqlString += s.join("LIMIT", "", []string{s.statement.limit}, "", "")
		}
	}

	if s.class == UPDATE {
		sqlString += s.join("UPDATE", "", s.statement.table, "", "")

		sqlString += s.join("JOIN", "", s.statement.join, "\nJOIN ", "")
		sqlString += s.join("INNER JOIN", "", s.statement.innerJoin, "\nINNER JOIN ", "")
		sqlString += s.join("OUTER JOIN", "", s.statement.outerJoin, "\nOUTER JOIN ", "")
		sqlString += s.join("LEFT OUTER JOIN", "", s.statement.leftOuterJoin, "\nLEFT OUTER JOIN ", "")
		sqlString += s.join("RIGHT OUTER JOIN", "", s.statement.rightOuterJoin, "\nRIGHT OUTER JOIN ", "")

		sqlString += s.join("SET", "", s.statement.set, ", ", "")
		sqlString += s.join("WHERE", "(", s.statement.where, " AND ", ")")

		if s.statement.offset != "" {
			sqlString += s.join("OFFSET", "", []string{s.statement.offset}, "", "")
		}
		if s.statement.limit != "" {
			sqlString += s.join("LIMIT", "", []string{s.statement.limit}, "", "")
		}
	}

	if s.class == DELETE {
		sqlString += s.join("DELETE FROM", "", s.statement.table, "", "")
		sqlString += s.join("WHERE", "(", s.statement.where, " AND ", ")")

		if s.statement.offset != "" {
			sqlString += s.join("OFFSET", "", []string{s.statement.offset}, "", "")
		}
		if s.statement.limit != "" {
			sqlString += s.join("LIMIT", "", []string{s.statement.limit}, "", "")
		}
	}

	if s.class == INSERT_INTO {
		sqlString += s.join("INSERT INTO", "", s.statement.table, "", "")
		sqlString += s.join("", "(", s.statement.columns, ", ", ")")
		sqlString += s.join("VALUES", "(", s.statement.values, ", ", ")")
	}

	if s.class == CREATE_TABLE {
		keyword := "CREATE TABLE"
		if s.statement.ifNotExists {
			keyword += " IF NOT EXISTS"
		}

		sqlString += s.join(keyword, "", s.statement.table, "", "")
		sqlString += s.join("", "(", s.statement.columns, ", ", ")")
	}

	if s.class == CREATE_SCHEMA {
		keyword := "CREATE SCHEMA"
		if s.statement.ifNotExists {
			keyword += " IF NOT EXISTS"
		}

		sqlString += s.join(keyword, "", []string{s.statement.schema}, "", "")
	}

	return strings.Trim(sqlString, "\n")
}

func (s *sqlBuilder) Params() []any {
	result := []any{}
	sqlString := s.String()
	matches := s.paramRegexp.FindAllString(sqlString, -1)

	for _, match := range matches {
		result = append(result, s.params[match])
	}

	return result
}
