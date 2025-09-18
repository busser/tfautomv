package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type jsonRule struct {
	baseRule
}

func parseJSONRule(s string) (*jsonRule, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, errors.New("syntax error")
	}

	r := jsonRule{
		baseRule: baseRule{
			resourceType: parts[0],
			attribute:    parts[1],
		},
	}

	return &r, nil
}

func (r jsonRule) String() string {
	return fmt.Sprintf("%s:%s:%s", RuleTypeJSON, r.resourceType, r.attribute)
}

func (r *jsonRule) Equates(a, b interface{}) bool {
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

	var aJSON, bJSON interface{}
	err := json.Unmarshal([]byte(aStr), &aJSON)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(bStr), &bJSON)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(aJSON, bJSON)
}
