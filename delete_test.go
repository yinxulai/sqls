package sqls

import "testing"

func TestSimpleDeleteStatement(t *testing.T) {
	s:= DELETE_FROM("PERSON")
	s.WHERE("P.ID = A.ID")
	s.WHERE("P.FIRST_NAME like ?")
	s.OR()
	s.WHERE("P.LAST_NAME like ?")
	s.OFFSET("10")
	s.LIMIT("20")

	result := s.String()
	expected := "DELETE FROM PERSON\nWHERE (P.ID = A.ID AND P.FIRST_NAME like ?) OR (P.LAST_NAME like ?)\nOFFSET 10\nLIMIT 20"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
