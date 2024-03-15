package sqls

import (
	"testing"
)

func TestCreateTableStatement(t *testing.T) {
	s:= CREATE_TABLE("PERSON")
	s.IF_NOT_EXISTS()
	s.COLUMN("id INT PRIMARY KEY")
	s.COLUMN("username VARCHAR(50) NOT NULL")
	s.COLUMN("email VARCHAR(100) UNIQUE")
	s.COLUMN("age INT")
	s.COLUMN("created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")

	result := s.String()
	expected := "CREATE TABLE IF NOT EXISTS PERSON (id INT PRIMARY KEY, username VARCHAR(50) NOT NULL, email VARCHAR(100) UNIQUE, age INT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"

	if result != expected {
		t.Errorf("Case1() 返回值为 %s，期望值为 %s", result, expected)
	}
}
