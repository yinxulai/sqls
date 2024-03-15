package sqls

import "testing"

func TestTruncateTableStatement(t *testing.T) {
	s := TRUNCATE_TABLE("PERSON")
	result := s.String()
	expected := "TRUNCATE TABLE PERSON"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
