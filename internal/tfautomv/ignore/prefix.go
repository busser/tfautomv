package ignore

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type prefixRule struct {
	baseRule
	prefix string
}

func parsePrefixRule(s string) (*prefixRule, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return nil, errors.New("syntax error")
	}

	r := prefixRule{
		baseRule: baseRule{
			resourceType: parts[0],
			attribute:    parts[1],
		},
		prefix: parts[2],
	}

	return &r, nil
}

func (r prefixRule) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", RuleTypePrefix, r.resourceType, r.attribute, r.prefix)
}

func (r *prefixRule) Equates(a, b interface{}) bool {
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

	return strings.TrimPrefix(aStr, r.prefix) == strings.TrimPrefix(bStr, r.prefix)
}
