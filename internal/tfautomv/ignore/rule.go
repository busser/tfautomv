package ignore

import (
	"errors"
	"fmt"
	"strings"
)

// A RuleType identifies a rule's logic.
type RuleType string

const (
	// RuleTypeEverything ignores all differences between two attributes'
	// values.
	RuleTypeEverything RuleType = "everything"

	// RuleTypePrefix ignores a given prefix when comparing attribute values.
	RuleTypePrefix RuleType = "prefix"

	// RuleTypeWhitespace ignores differences in whitespace between two
	// attributes' values. Whitespace is as defined by unicode.IsSpace.
	RuleTypeWhitespace RuleType = "whitespace"
)

// A Rule allows tfautomv to equate certains attribute values that would
// normally be considered different.
type Rule interface {
	// Each rule has a unique string representation. This representation is how
	// users provides rules to tfautomv.
	fmt.Stringer

	// AppliesTo returns wether the Rule applies to the given resource type and
	// attribute.
	AppliesTo(resourceType, attribute string) bool

	// Equates checks whether two values match, according to the Rule's logic.
	Equates(a, b interface{}) bool
}

// ParseRule converts a string into a Rule.
func ParseRule(s string) (Rule, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) < 2 {
		return nil, errors.New("invalid syntax")
	}

	ruleType := RuleType(parts[0])

	switch ruleType {
	case RuleTypeEverything:
		return parseEverythingRule(parts[1])
	case RuleTypePrefix:
		return parsePrefixRule(parts[1])
	case RuleTypeWhitespace:
		return parseWhitespaceRule(parts[1])
	default:
		return nil, fmt.Errorf("unknown rule type %q", ruleType)
	}
}

// MustParseRule converts a string into a Rule. MustParseRule panics if the
// string is not a valid rule.
func MustParseRule(s string) Rule {
	r, err := ParseRule(s)
	if err != nil {
		panic(fmt.Sprintf("MustParseRule(): %v", err))
	}
	return r
}
