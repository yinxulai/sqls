package sqls

import "testing"

func TestSimpleUpdateStatement(t *testing.T) {
	s := UPDATE("PERSON")
	s.INNER_JOIN("DEPARTMENT D on D.ID = P.DEPARTMENT_ID")
	s.INNER_JOIN("COMPANY C on D.COMPANY_ID = C.ID")
	s.SET("a", "1")
	s.SET("b", "2")
	s.WHERE("P.ID = A.ID")
	s.WHERE("P.FIRST_NAME like ?")
	s.OR()
	s.WHERE("P.LAST_NAME like ?")

	result := s.String()
	expected := "UPDATE PERSON\nINNER JOIN DEPARTMENT D on D.ID = P.DEPARTMENT_ID\nINNER JOIN COMPANY C on D.COMPANY_ID = C.ID\nSET a=1, b=2\nWHERE (P.ID = A.ID AND P.FIRST_NAME like ?) OR (P.LAST_NAME like ?)"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
