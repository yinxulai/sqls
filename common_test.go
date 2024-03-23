package sqls

import (
	"slices"
	"testing"
)

func TestCommonParamRegexp(t *testing.T) {
	caseText := "$123"
	result := paramRegexp.FindAllString(caseText, -1)

	expected := []string{"$123"}

	if slices.Compare[[]string](result, expected) != 0 {
		t.Errorf("FindAllString(%s) 返回值为 %v，期望值为 %v", caseText, result, expected)
	}
}

func TestSimplePostgresqlParam(t *testing.T) {
	builder := newSqlBuilder()

	caseText := builder.Param(1) + " " + builder.Param(2) + " " + builder.Param(3)
	result := paramRegexp.FindAllString(caseText, -1)

	expected := []string{"$1", "$2", "$3"}

	if slices.Compare[[]string](result, expected) != 0 {
		t.Errorf("FindAllString(%s) 返回值为 %v，期望值为 %v", caseText, result, expected)
	}

	params := builder.Params(caseText)
	paramSlice := []int{params[0].(int), params[1].(int), params[2].(int)}
	expectedParams := []int{1, 2, 3}
	if slices.Compare[[]int](paramSlice, expectedParams) != 0 {
		t.Errorf("Params(%s) 返回值为 %v，期望值为 %v", caseText, result, expected)
	}
}
