package engine

import "testing"

func TestCompareAll(t *testing.T) {
	tests := []struct {
		name   string
		create []Resource
		delete []Resource

		wantComparisonCount int
	}{
		{
			// There are two cases where we may want to compare a resource to
			// one with the exact same address. In each case, we have a reason
			// not to perform the comparison.
			//
			// 1. The resource is being created and deleted in the same plan,
			//    not because of refactoring but because of a desired change.
			//    This violates tfautomv's assumption that the plan was empty
			//    before refactoring.
			// 2. The resource was moved to another address and a second
			//    resource was moved to the original address. Such an operation
			//    cannot be done in a single sweep with Terraform. Since the
			//    two moves would need to be done in a specific order, which we
			//    cannot guess, this is not a valid refactoring. Instead,
			//    tfautomv should be run between each move.
			name: "same module same type same address",
			create: []Resource{
				dummyResource("same_module", "same_type", "same_address"),
			},
			delete: []Resource{
				dummyResource("same_module", "same_type", "same_address"),
			},
			wantComparisonCount: 0,
		},

		{
			// The very base case of a refactoring: a resource is moved to a
			// different address within the same module.
			name: "same module same type different address",
			create: []Resource{
				dummyResource("same_module", "same_type", "this_address"),
			},
			delete: []Resource{
				dummyResource("same_module", "same_type", "that_address"),
			},
			wantComparisonCount: 1,
		},

		{
			// Since the two resources have a different type, one cannot be
			// moved to the other's address. This is not a valid refactoring.
			name: "same module different type different address",
			create: []Resource{
				dummyResource("same_module", "this_type", "this_address"),
			},
			delete: []Resource{
				dummyResource("same_module", "that_type", "that_address"),
			},
			wantComparisonCount: 0,
		},

		{
			// A resource is moved to a different module. This is a valid
			// refactoring.
			name: "different module same type same address",
			create: []Resource{
				dummyResource("this_module", "same_type", "same_address"),
			},
			delete: []Resource{
				dummyResource("that_module", "same_type", "same_address"),
			},
			wantComparisonCount: 1,
		},

		{
			// A resource is moved to a different module and a different
			// address. This is a valid refactoring.
			name: "different module same type different address",
			create: []Resource{
				dummyResource("this_module", "same_type", "this_address"),
			},
			delete: []Resource{
				dummyResource("that_module", "same_type", "that_address"),
			},
			wantComparisonCount: 1,
		},

		{
			// The two resources are entirely unrelated.
			name: "different module different type different address",
			create: []Resource{
				dummyResource("this_module", "this_type", "this_address"),
			},
			delete: []Resource{
				dummyResource("that_module", "that_type", "that_address"),
			},
			wantComparisonCount: 0,
		},

		{
			// In this case, we want to compare the three resources to create
			// to the three resources to delete. However, we don't want to
			// compare the resources to create to each other, nor do we want to
			// compare the resources to delete to each other.
			name: "many resources with same type",
			create: []Resource{
				dummyResource("this_module", "same_type", "this_address"),
				dummyResource("this_module", "same_type", "that_address"),
				dummyResource("that_module", "same_type", "other_address"),
			},
			delete: []Resource{
				dummyResource("that_module", "same_type", "this_address"),
				dummyResource("that_module", "same_type", "that_address"),
				dummyResource("this_module", "same_type", "other_address"),
			},
			wantComparisonCount: 9,
		},

		{
			// In this case, we want to compare the resources to create to the
			// resources to delete of the same type.
			name: "many resources with few different types",
			create: []Resource{
				dummyResource("this_module", "this_type", "this_address"),
				dummyResource("this_module", "this_type", "that_address"),
				dummyResource("that_module", "that_type", "other_address"),
			},
			delete: []Resource{
				dummyResource("that_module", "that_type", "this_address"),
				dummyResource("that_module", "that_type", "that_address"),
				dummyResource("this_module", "this_type", "other_address"),
			},
			wantComparisonCount: 4,
		},

		{
			// In this case, there is a single resource to create and a single
			// resource to delete for each type. We want to compare these pairs
			// of resources. We do not want to compare the resources with
			// different types.
			name: "many resources with many different types",
			create: []Resource{
				dummyResource("this_module", "this_type", "this_address"),
				dummyResource("this_module", "that_type", "that_address"),
				dummyResource("that_module", "other_type", "other_address"),
			},
			delete: []Resource{
				dummyResource("that_module", "this_type", "this_address"),
				dummyResource("that_module", "that_type", "that_address"),
				dummyResource("this_module", "other_type", "other_address"),
			},
			wantComparisonCount: 3,
		},

		{
			// In this case, there are no resources to create. There is no
			// reason to compare any resources.
			name:   "no resources to create",
			create: []Resource{},
			delete: []Resource{
				dummyResource("this_modue", "this_type", "this_address"),
				dummyResource("this_module", "that_type", "that_address"),
				dummyResource("that_module", "that_type", "other_address"),
			},
			wantComparisonCount: 0,
		},

		{
			// In this case, there are no resources to delete. There is no
			// reason to compare any resources.
			name: "no resources to delete",
			create: []Resource{
				dummyResource("this_modue", "this_type", "this_address"),
				dummyResource("this_module", "that_type", "that_address"),
				dummyResource("that_module", "that_type", "other_address"),
			},
			delete:              []Resource{},
			wantComparisonCount: 0,
		},

		{
			// In this case, the resources to create have a different type than
			// the resources to delete. There is no reason to compare any
			// resources.
			name: "many resources with all different types",
			create: []Resource{
				dummyResource("this_module", "this_type", "this_address"),
				dummyResource("this_module", "that_type", "that_address"),
				dummyResource("that_module", "other_type", "other_address"),
			},
			delete: []Resource{
				dummyResource("this_module", "another_type", "this_address"),
				dummyResource("that_module", "yet_another_type", "that_address"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := Plan{
				ToCreate: tt.create,
				ToDelete: tt.delete,
			}

			comparisons := CompareAll(plan, nil)
			if got, want := len(comparisons), tt.wantComparisonCount; got != want {
				t.Errorf("got %d comparisons, want %d", got, want)
			}
		})
	}
}

func dummyResource(moduleID, typ, address string) Resource {
	if moduleID == "" {
		moduleID = "dummy_module"
	}
	if typ == "" {
		typ = "dummy_type"
	}
	if address == "" {
		address = "dummy_address"
	}

	return Resource{
		ModuleID:   moduleID,
		Type:       typ,
		Address:    address,
		Attributes: dummyAttributes,
	}
}

var dummyAttributes = map[string]any{
	"a.b.c": "foo",
	"d.e.0": "bar",
	"d.e.1": "baz",
	"d.e.#": 2,
	"f":     false,
}
