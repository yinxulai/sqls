package sqls

import (
	"slices"
	"testing"
)

func TestCommonParamRegexp(t *testing.T) {
	caseText := "$123"
	result := paramRegexp.FindAllStringSubmatch(caseText, -1)

	expected := []string{"$123"}

	if slices.Compare[[]string](result[0], expected) != 0 {
		t.Errorf("FindAllString(%s) 返回值为 %v，期望值为 %v", caseText, result[0], expected)
	}
}
