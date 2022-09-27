package ignore

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type whitespaceRule struct {
	baseRule
}

func parseWhitespaceRule(s string) (*whitespaceRule, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, errors.New("syntax error")
	}

	r := whitespaceRule{
		baseRule: baseRule{
			resourceType: parts[0],
			attribute:    parts[1],
		},
	}

	return &r, nil
}

func (r whitespaceRule) String() string {
	return fmt.Sprintf("%s:%s:%s", RuleTypeWhitespace, r.resourceType, r.attribute)
}

func (r *whitespaceRule) Equates(a, b interface{}) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != bVal.Kind() {
		return false
	}
	kind := aVal.Kind()

	var aStr, bStr string
	if kind == reflect.String {
		aStr = aVal.String()
		bStr = bVal.String()
	} else {
		aStr = fmt.Sprint(a)
		bStr = fmt.Sprint(b)
	}

	return withoutWhitespace(aStr) == withoutWhitespace(bStr)
}

func withoutWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
