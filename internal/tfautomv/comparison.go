package tfautomv

// A Comparison contains a list matching and mismatching attributes between two
// resources.
type Comparison struct {
	// The resources that were compared.
	Created   *Resource
	Destroyed *Resource

	// Attributes that are set in Created and have the same value in Destroyed.
	MatchingAttributes []string

	// Attributes that are set in Destroyed and are not set or do not have the
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
// Attributes set in destroyed but not in created are ignored. We assume they
// are set by the Terraform provider, the cloud provider, or an external actor.
func Compare(created, destroyed *Resource) Comparison {
	comp := Comparison{
		Created:   created,
		Destroyed: destroyed,
	}

	for attr, createdVal := range created.Attributes {
		if createdVal == nil {
			continue
		}

		destroyedVal, isSet := destroyed.Attributes[attr]
		if isSet && createdVal == destroyedVal {
			comp.MatchingAttributes = append(comp.MatchingAttributes, attr)
		} else {
			comp.MismatchingAttributes = append(comp.MismatchingAttributes, attr)
		}
	}

	return comp
}

func (c *Comparison) IsMatch() bool {
	// Resources match if none of their attributes mismatch.
	return len(c.MismatchingAttributes) == 0
}
