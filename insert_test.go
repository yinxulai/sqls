package sqls

import "testing"

func TestSimpleInsertStatement(t *testing.T) {
	s:= INSERT_INTO("PERSON")
	s.VALUES("a", s.Param(1))
	s.VALUES("b", s.Param(1))
	s.VALUES("c", s.Param(1))
	s.ON_CONFLICT("a, b")
	s.ON_CONFLICT("c")
	s.DO_UPDATE_SET("a","1")
	s.DO_UPDATE_SET("b","2")

	result := s.String()
	expected := "INSERT INTO PERSON (a, b, c) VALUES ($1, $2, $3) ON CONFLICT (a, b, c) DO UPDATE SET a=1, b=2"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
