package tfautomv

import (
	"sort"
	"testing"

	"github.com/padok-team/tfautomv/internal/slices"
	"github.com/padok-team/tfautomv/internal/terraform"
)

type dummyResource struct {
	address string
	typ     string
}

func dummyPlan(t *testing.T, created, destroyed []dummyResource) *terraform.Plan {
	t.Helper()

	var plan terraform.Plan

	for _, r := range created {
		plan.ResourceChanges = append(plan.ResourceChanges, terraform.ResourceChange{
			Address: r.address,
			Type:    r.typ,
			Change: terraform.Change{
				Actions: []string{terraform.CreateAction},
			},
		})
	}
	for _, r := range destroyed {
		plan.ResourceChanges = append(plan.ResourceChanges, terraform.ResourceChange{
			Address: r.address,
			Type:    r.typ,
			Change: terraform.Change{
				Actions: []string{terraform.DeleteAction},
			},
		})
	}

	return &plan
}

func checkResourcesByTypeMatch(t *testing.T, operation string, actual map[string][]*Resource, want map[string][]string) {
	t.Helper()

	if len(actual) != len(want) {
		t.Fatalf("%s: got %d types, want %d",
			operation, len(actual), len(want))
	}

	for typ, resources := range actual {
		var actualAddresses []string
		for _, r := range resources {
			actualAddresses = append(actualAddresses, r.Address)
		}
		wantAddresses := want[typ]

		sort.Strings(actualAddresses)
		sort.Strings(wantAddresses)

		if !slices.Equal(actualAddresses, wantAddresses) {
			t.Errorf("%s: got %q, want %q",
				operation, actualAddresses, wantAddresses)
		}
	}
}

func TestAnalysisFromPlan(t *testing.T) {
	tt := []struct {
		createdResources    []dummyResource
		destroyedResources  []dummyResource
		wantCreatedByType   map[string][]string
		wantDestroyedByType map[string][]string
		wantComparisons     map[string][]string
	}{
		{
			createdResources: []dummyResource{
				{"a", "type-0"},
				{"b", "type-0"},
				{"c", "type-1"},
				{"d", "type-1"},
				{"e", "type-2"},
				{"f", "type-3"},
			},
			destroyedResources: []dummyResource{
				{"g", "type-0"},
				{"h", "type-1"},
				{"i", "type-1"},
				{"j", "type-2"},
				{"k", "type-0"},
				{"l", "type-4"},
			},
			wantCreatedByType: map[string][]string{
				"type-0": {"a", "b"},
				"type-1": {"c", "d"},
				"type-2": {"e"},
				"type-3": {"f"},
			},
			wantDestroyedByType: map[string][]string{
				"type-0": {"g", "k"},
				"type-1": {"h", "i"},
				"type-2": {"j"},
				"type-4": {"l"},
			},
			wantComparisons: map[string][]string{
				"a": {"g", "k"},
				"b": {"g", "k"},
				"c": {"h", "i"},
				"d": {"h", "i"},
				"e": {"j"},
				"g": {"a", "b"},
				"h": {"c", "d"},
				"i": {"c", "d"},
				"j": {"e"},
				"k": {"a", "b"},
			},
		},
	}

	for _, tc := range tt {
		plan := dummyPlan(t, tc.createdResources, tc.destroyedResources)

		actual, err := AnalysisFromPlan(plan, nil)
		if err != nil {
			t.Fatalf("AnalysisFromPlan(): unexpected error: %v", err)
		}

		checkResourcesByTypeMatch(t, "created", actual.CreatedByType, tc.wantCreatedByType)
		checkResourcesByTypeMatch(t, "destroyed", actual.DestroyedByType, tc.wantDestroyedByType)

		if len(actual.Comparisons) != len(tc.wantComparisons) {
			t.Fatalf("got %d resources with comparisons, want %d",
				len(actual.Comparisons), len(tc.wantComparisons))
		}

		for res, comps := range actual.Comparisons {
			wantAddresses, ok := tc.wantComparisons[res.Address]
			if !ok {
				t.Errorf("did not expect resource with address %q to have comparisons",
					res.Address)
			}

			var gotAddresses []string
			for _, c := range comps {
				addr := c.Created.Address
				if addr == res.Address {
					addr = c.Destroyed.Address
				}
				gotAddresses = append(gotAddresses, addr)
			}

			sort.Strings(gotAddresses)
			sort.Strings(wantAddresses)

			if !slices.Equal(gotAddresses, wantAddresses) {
				t.Errorf("got %q compared to %q, want %q",
					res.Address, gotAddresses, wantAddresses)
			}
		}
	}
}
