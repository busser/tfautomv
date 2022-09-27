package ignore

import "testing"

func TestEverythingRuleAppliesTo(t *testing.T) {
	rule := everythingRule{
		baseRule{
			resourceType: "my_resource",
			attribute:    "my_attr",
		},
	}

	tt := []struct {
		resourceType string
		attribute    string
		want         bool
	}{
		{
			resourceType: "my_resource",
			attribute:    "my_attr",
			want:         true,
		},
		{
			resourceType: "not_my_resource",
			attribute:    "my_attr",
			want:         false,
		},
		{
			resourceType: "my_resource",
			attribute:    "not_my_attr",
			want:         false,
		},
		{
			resourceType: "not_my_resource",
			attribute:    "not_my_attr",
			want:         false,
		},
	}

	for _, tc := range tt {
		actual := rule.AppliesTo(tc.resourceType, tc.attribute)
		if actual != tc.want {
			t.Errorf("AppliesTo(%q, %q) = %t, want %t", tc.resourceType, tc.attribute, actual, tc.want)
		}
	}
}

func TestEverythingRuleEquates(t *testing.T) {
	rule := everythingRule{
		baseRule{
			resourceType: "my_resource",
			attribute:    "my_attr",
		},
	}

	values := []interface{}{
		"foo",
		"\tfoo bar\n",
		123,
		456,
		"123",
		true,
		false,
		"false",
	}

	for _, a := range values {
		for _, b := range values {
			if !rule.Equates(a, b) {
				t.Errorf("Equates(%q, %q) = false, want true", a, b)
			}
		}
	}
}
