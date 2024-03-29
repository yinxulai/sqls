package sqls

import (
	"testing"
)

func Case1() string {
	s := SELECT("P.ID, P.USERNAME, P.PASSWORD, P.FULL_NAME")
	s.SELECT("P.LAST_NAME, P.CREATED_ON, P.UPDATED_ON")
	s.FROM("PERSON P")
	s.FROM("ACCOUNT A")
	s.INNER_JOIN("DEPARTMENT D on D.ID = P.DEPARTMENT_ID")
	s.INNER_JOIN("COMPANY C on D.COMPANY_ID = C.ID")
	s.WHERE("P.ID = A.ID")
	s.WHERE("P.FIRST_NAME like ?")
	s.OR()
	s.WHERE("P.LAST_NAME like ?")
	s.GROUP_BY("P.ID")
	s.HAVING("P.LAST_NAME like ?")
	s.OR()
	s.HAVING("P.FIRST_NAME like ?")
	s.ORDER_BY("P.ID")
	s.ORDER_BY("P.FULL_NAME")
	s.OFFSET("10")
	s.LIMIT("20")

	return s.String()
}

func Case2(id *string, firstName *string, lastName *string) string {
	s := SELECT("P.ID, P.USERNAME, P.PASSWORD, P.FIRST_NAME, P.LAST_NAME")
	s.FROM("PERSON P")
	if id != nil {
		s.WHERE("P.ID like #id#")
	}
	if firstName != nil {
		s.WHERE("P.FIRST_NAME like #firstName#")
	}
	if lastName != nil {
		s.WHERE("P.LAST_NAME like #lastName#")
	}
	s.ORDER_BY("P.LAST_NAME")
	return s.String()
}

func TestSimpleSelectStatement(t *testing.T) {
	var a = "a"
	var b = "b"
	var c = "c"

	result := Case2(&a, &b, &c)
	expected := "SELECT P.ID, P.USERNAME, P.PASSWORD, P.FIRST_NAME, P.LAST_NAME FROM PERSON P WHERE (P.ID like #id# AND P.FIRST_NAME like #firstName# AND P.LAST_NAME like #lastName#) ORDER BY P.LAST_NAME"

	if result != expected {
		t.Errorf("Case2(&a, &b, &c) 返回值为 %s，期望值为 %s", result, expected)
	}
}

func TestSimpleSelectStatementMissingFirstParam(t *testing.T) {
	var b = "b"
	var c = "c"

	result := Case2(nil, &b, &c)
	expected := "SELECT P.ID, P.USERNAME, P.PASSWORD, P.FIRST_NAME, P.LAST_NAME FROM PERSON P WHERE (P.FIRST_NAME like #firstName# AND P.LAST_NAME like #lastName#) ORDER BY P.LAST_NAME"

	if result != expected {
		t.Errorf("Case2(nil, &b, &c) 返回值为 %s，期望值为 %s", result, expected)
	}
}

func TestSimpleSelectStatementMissingFirstTwoParams(t *testing.T) {
	var c = "c"

	result := Case2(nil, nil, &c)
	expected := "SELECT P.ID, P.USERNAME, P.PASSWORD, P.FIRST_NAME, P.LAST_NAME FROM PERSON P WHERE (P.LAST_NAME like #lastName#) ORDER BY P.LAST_NAME"

	if result != expected {
		t.Errorf("Case2(nil, nil, &c) 返回值为 %s，期望值为 %s", result, expected)
	}
}

func TestSimpleSelectStatementMissingAllParams(t *testing.T) {
	result := Case2(nil, nil, nil)
	expected := "SELECT P.ID, P.USERNAME, P.PASSWORD, P.FIRST_NAME, P.LAST_NAME FROM PERSON P ORDER BY P.LAST_NAME"

	if result != expected {
		t.Errorf("Case2(nil, nil, nil) 返回值为 %s，期望值为 %s", result, expected)
	}
}

func TestComplexSelectStatement(t *testing.T) {
	result := Case1()
	expected := "SELECT P.ID, P.USERNAME, P.PASSWORD, P.FULL_NAME, P.LAST_NAME, P.CREATED_ON, P.UPDATED_ON FROM PERSON P, ACCOUNT A INNER JOIN DEPARTMENT D on D.ID = P.DEPARTMENT_ID INNER JOIN COMPANY C on D.COMPANY_ID = C.ID WHERE (P.ID = A.ID AND P.FIRST_NAME like ?) OR (P.LAST_NAME like ?) GROUP BY P.ID HAVING (P.LAST_NAME like ?) OR (P.FIRST_NAME like ?) ORDER BY P.ID, P.FULL_NAME OFFSET 10 LIMIT 20"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
