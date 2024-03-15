package sqls

import (
	"fmt"
	"regexp"
	"slices"
)

type sqlBuilder struct {
	params map[string]any
}

var paramRegexp = regexp.MustCompile(`\$\d+`)

func newSqlBuilder() *sqlBuilder {
	builder := &sqlBuilder{}
	builder.params = make(map[string]any)
	return builder
}

var orString = ") OR ("
var andString = ") AND ("
var condSlice = []string{orString, andString}

// 拼接 sql 的工具方法
func (s *sqlBuilder) join(keyword string, open string, fields []string, sep string, close string) string {
	if len(fields) == 0 {
		return ""
	}

	var body = ""
	for index, field := range fields {
		body += field
		if index != len(fields)-1 {
			nextField := fields[index+1]
			if !slices.Contains[[]string](condSlice, field) && !slices.Contains[[]string](condSlice, nextField) {
				body += sep
			}
		}
	}

	if keyword != "" {
		return keyword + " " + open + body + close + " "
	}

	return open + body + close + " "
}

func (s *sqlBuilder) Param(v any) string {
	paramsIndex := len(s.params) + 1
	key := fmt.Sprintf("$%d", paramsIndex)
	s.params[key] = v
	return key
}

func (s *sqlBuilder) Params(sql string) []any {
	result := []any{}
	matches := paramRegexp.FindAllString(sql, -1)

	for _, match := range matches {
		result = append(result, s.params[match])
	}

	return result
}
