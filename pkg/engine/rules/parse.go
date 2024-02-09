package rules

import (
	"errors"
	"fmt"
	"strings"

	"github.com/busser/tfautomv/pkg/engine"
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

// Parse converts a string into a Rule.
func Parse(s string) (engine.Rule, error) {
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

// MustParse converts a string into a Rule. MustParse panics if the
// string is not a valid rule.
func MustParse(s string) engine.Rule {
	r, err := Parse(s)
	if err != nil {
		panic(fmt.Sprintf("MustParseRule(): %v", err))
	}
	return r
}
