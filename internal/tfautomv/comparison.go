package tfautomv

import "github.com/padok-team/tfautomv/internal/tfautomv/ignore"

// A Comparison contains a list matching and mismatching attributes between two
// resources.
type Comparison struct {
	// The resources that were compared.
	Created   *Resource
	Destroyed *Resource

	// Attributes that are set in Created and have the same value in Destroyed.
	MatchingAttributes []string

	// Attributes that are set in Created and are not set or do not have the
	// same value in Destroyed, but those differences are ignored because of a
	// rule.
	IgnoredAttributes []string

	// Attributes that are set in Created and are not set or do not have the
	// same value in Destroyed.
	MismatchingAttributes []string
}

// Compare finds which attributes match between two resources: one planned for
// creation and the other planned for destruction.
//
// An attribute matches if it is set in created and is set and has the same
// value in destroyed. An attribute mismatches if it is set in created and is
// not set or has a different value in destroyed.
//
// Specific rules make it so certain differences are ignored and will not
// trigger a mismatch.
//
// Attributes set in destroyed but not in created are ignored. We assume they
// are set by the Terraform provider, the cloud provider, or an external actor.
func Compare(created, destroyed *Resource, rules []ignore.Rule) Comparison {
	comp := Comparison{
		Created:   created,
		Destroyed: destroyed,
	}

	for attr, createdVal := range created.Attributes {
		if createdVal == nil {
			continue
		}

		destroyedVal, isSet := destroyed.Attributes[attr]

		// Match: both values are identical.
		if isSet && createdVal == destroyedVal {
			comp.MatchingAttributes = append(comp.MatchingAttributes, attr)
			continue
		}

		// Ignored mismatch: differences between values are ignored.
		var ignored bool
		for _, r := range rules {
			if !r.AppliesTo(created.Type, attr) {
				continue
			}
			if r.Equates(createdVal, destroyedVal) {
				ignored = true
				comp.IgnoredAttributes = append(comp.IgnoredAttributes, attr)
				break
			}
		}
		if ignored {
			continue
		}

		// Mismatch: values are different and no rules ignore those differences.
		comp.MismatchingAttributes = append(comp.MismatchingAttributes, attr)
	}

	return comp
}

func (c *Comparison) IsMatch() bool {
	// Resources match if none of their attributes mismatch.
	return len(c.MismatchingAttributes) == 0
}
