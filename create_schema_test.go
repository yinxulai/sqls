package sqls

import "testing"

func TestCreateSchemaStatement(t *testing.T) {
	s := CREATE_SCHEMA("PERSON").IF_NOT_EXISTS()
	result := s.String()
	expected := "CREATE SCHEMA IF NOT EXISTS PERSON"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
