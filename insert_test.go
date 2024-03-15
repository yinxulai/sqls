package sqls

import "testing"

func TestSimpleInsertStatement(t *testing.T) {
	s:= INSERT_INTO("PERSON")
	s.VALUES("a", s.Param(1))
	s.VALUES("b", s.Param(1))
	s.VALUES("c", s.Param(1))

	result := s.String()
	expected := "INSERT INTO PERSON (a, b, c) VALUES ($1, $2, $3)"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
