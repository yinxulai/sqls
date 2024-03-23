package sqls

import (
	"fmt"
	"strings"
)

type constraint struct {
	name    string
	options string
}

type alterTableBuilder struct {
	builder       *sqlBuilder
	table         string
	addConstraint []constraint
}

func newAlterTableBuilder() *alterTableBuilder {
	builder := &alterTableBuilder{}
	builder.builder = newSqlBuilder()
	builder.addConstraint = make([]constraint, 0)
	return builder
}

func ALTER_TABLE(table string) *alterTableBuilder {
	s := newAlterTableBuilder()
	s.table = table
	return s
}

func (s *alterTableBuilder) ADD_CONSTRAINT(name string, options string) *alterTableBuilder {
	s.addConstraint = append(s.addConstraint, constraint{name: name, options: options})
	return s
}

func (s *alterTableBuilder) Param(v any) string {
	return s.builder.Param(v)
}

func (s *alterTableBuilder) Params() []any {
	return s.builder.Params(s.String())
}

func (s *alterTableBuilder) String() string {
	var sqlString string
	sqlString += s.builder.join("ALTER TABLE", "", []string{s.table}, "", "")

	if len(s.addConstraint) > 0 {
		var addConstraint []string
		for _, constraint := range s.addConstraint {
			addConstraint = append(addConstraint, fmt.Sprintf("ADD CONSTRAINT %s %s", constraint.name, constraint.options))
		}
		sqlString += s.builder.join("", "", addConstraint, ", ", "")
	}

	return strings.Trim(sqlString, " ")
}
