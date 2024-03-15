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
