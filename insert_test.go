package sqls

import "testing"

func TestSimpleInsertStatement(t *testing.T) {
	s:= INSERT_INTO("PERSON")
	s.VALUES("a", "a")
	s.VALUES("b", "a")
	s.VALUES("c", "a")

	result := s.String()
	expected := "INSERT INTO PERSON\n(a, b, c)\nVALUES (a, a, a)"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
