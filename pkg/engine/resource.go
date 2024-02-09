package engine

import (
	"fmt"
	"sort"
)

// A Resource represents a Terraform resource. Whether Terraform plans
// to create or delete the resource does not matter; the resource is modeled in
// the same way in both cases.
type Resource struct {
	// The module where the resource is defined.
	ModuleID string

	// The resource's type.
	Type string

	// The resource's address within the module's state. Terraform uses the
	// address to map a module's source code to resources it manages.
	Address string

	// The resource's known attributes. In Terraform's state, attributes are
	// complex objects. We flatten them to make them easier to work with.
	//
	// For example, the following Terraform resource:
	//
	//	resource "aws_instance" "web" {
	//	  ami           = "ami-a1b2c3d4"
	//	  instance_type = "t2.micro"
	//
	//	  tags = {
	//	    Name        = "HelloWorld"
	//	    Environment = "Production"
	//	  }
	//	}
	//
	// Would be represented as the following Attributes:
	//
	//	Attributes{
	//	  "ami":              "ami-a1b2c3d4",
	//	  "instance_type":    "t2.micro",
	//	  "tags.Name":        "HelloWorld",
	//	  "tags.Environment": "Production",
	//	}
	//
	// Note that the "tags" attribute is flattened into "tags.Name" and
	// "tags.Environment".
	//
	// A value in the flattened map is either a string, a number, a boolean, or
	// nil. The nil value is used to represent null values in Terraform.
	Attributes map[string]any
}

// A unique ID for the resource, for use as map keys. This ID is a concatenation
// of the module ID and the resource's address. Whether the resource is planned
// for creation or deletion is not taken into account.
func (r Resource) ID() string {
	return fmt.Sprintf("%s:%s", r.ModuleID, r.Address)
}

// A ResourceComparison represents a pair of Terraform resources of the same
// type: one that Terraform plans to create and another that Terraform plans to
// delete. By comparing these resources, we can determine whether we should move
// the existing resource's state to the new resource's address.
//
// An attribute is considered matching when both resources have the same value
// for that attribute.
// An attribute is considered mismatching when both resources have different
// values for that attribute.
// An attribute is considered ignored when both resources have different values
// for that attribute, but a rule says to ignore the difference between those
// values.
//
// A ResourceComparison is considered a match when no attributes are
// mismatching.
type ResourceComparison struct {
	// The resource Terraform plans to create.
	ToCreate Resource

	// The resource Terraform plans to delete.
	ToDelete Resource

	// Keys of attributes that have the same value in both resources.
	MatchingAttributes []string

	// Keys of attributes that have different values in both resources.
	MismatchingAttributes []string

	// Keys of attributes that would normally be mismatching, but where the user
	// provided a rule that says to ignore that particular difference.
	IgnoredAttributes []string
}

// IsMatch returns whether the two resources are a match.
func (rc ResourceComparison) IsMatch() bool {
	return len(rc.MismatchingAttributes) == 0
}

// CompareResources compares the attributes of two Terraform resources: one that
// Terraform plans to create and another that Terraform plans to delete.
func CompareResources(create, delete Resource, rules []Rule) ResourceComparison {
	var matching, mismatching, ignored []string

	for key, cValue := range create.Attributes {
		if cValue == nil {
			continue
		}

		dValue, isSet := delete.Attributes[key]

		if isSet && cValue == dValue {
			// Both values are identical: it's a match.
			matching = append(matching, key)
			continue
		}

		var ruleSaysToIgnore bool
		for _, r := range rules {
			if !r.AppliesTo(create.Type, key) {
				continue
			}

			if r.Equates(cValue, dValue) {
				ruleSaysToIgnore = true
				break
			}
		}

		if ruleSaysToIgnore {
			// A rule says to ignore the difference between the two values.
			ignored = append(ignored, key)
			continue
		}

		// The two values are different and no rule says to ignore the difference.
		mismatching = append(mismatching, key)
	}

	// We sort the keys so that the final diff is deterministic.
	sort.Strings(matching)
	sort.Strings(mismatching)
	sort.Strings(ignored)

	return ResourceComparison{
		ToCreate:              create,
		ToDelete:              delete,
		MatchingAttributes:    matching,
		MismatchingAttributes: mismatching,
		IgnoredAttributes:     ignored,
	}
}
