package ignore

import (
	"errors"
	"fmt"
	"strings"
)

type everythingRule struct {
	baseRule
}

func parseEverythingRule(s string) (*everythingRule, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, errors.New("syntax error")
	}

	r := everythingRule{
		baseRule: baseRule{
			resourceType: parts[0],
			attribute:    parts[1],
		},
	}

	return &r, nil
}

func (r everythingRule) String() string {
	return fmt.Sprintf("%s:%s:%s", RuleTypeEverything, r.resourceType, r.attribute)
}

func (r *everythingRule) Equates(a, b interface{}) bool {
	return true
}
