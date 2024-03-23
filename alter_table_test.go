package sqls

import (
	"testing"
)

func TestAlterTableAddConstraint(t *testing.T) {
	s := ALTER_TABLE("TABLE")
	s.ADD_CONSTRAINT("test", "UNIQUE (column1, column2)")
	s.ADD_CONSTRAINT("test2", "PRIMARY KEY (column1, column2)")

	result := s.String()
	expected := "ALTER TABLE TABLE ADD CONSTRAINT test UNIQUE (column1, column2), ADD CONSTRAINT test2 PRIMARY KEY (column1, column2)"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
