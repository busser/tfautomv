package ignore

import "testing"

func TestPrefixRuleAppliesTo(t *testing.T) {
	rule := prefixRule{
		baseRule{
			resourceType: "my_resource",
			attribute:    "my_attr",
		},
		"does-not-matter",
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

func TestPrefixRuleEquates(t *testing.T) {
	tt := []struct {
		valueA interface{}
		valueB interface{}
		prefix string
		want   bool
	}{
		{
			valueA: "foo",
			valueB: "foo",
			prefix: "any",
			want:   true,
		},
		{
			valueA: "foo",
			valueB: "foo",
			prefix: "",
			want:   true,
		},
		{
			valueA: "bc",
			valueB: "abc",
			prefix: "a",
			want:   true,
		},
		{
			valueA: "ac",
			valueB: "abc",
			prefix: "b",
			want:   false,
		},
		{
			valueA: "ab",
			valueB: "abc",
			prefix: "c",
			want:   false,
		},
		{
			valueA: "\tbar",
			valueB: "foo\tbar",
			prefix: "foo",
			want:   true,
		},
		{
			valueA: "qwertyuiop",
			valueB: "b/qwertyuiop",
			prefix: "b/",
			want:   true,
		},
		{
			valueA: "\tbar\n",
			valueB: "foo\tbar",
			prefix: "foo",
			want:   false,
		},
		{
			valueA: 123,
			valueB: 123,
			prefix: "any",
			want:   true,
		},
		{
			valueA: 123,
			valueB: 456,
			prefix: "any",
			want:   false,
		},
		{
			valueA: 123,
			valueB: "123",
			prefix: "any",
			want:   false,
		},
		{
			valueA: 123,
			valueB: "foo123",
			prefix: "foo",
			want:   false,
		},
		{
			valueA: false,
			valueB: false,
			prefix: "any",
			want:   true,
		},
		{
			valueA: false,
			valueB: "false",
			prefix: "any",
			want:   false,
		},
	}

	for _, tc := range tt {
		rule := prefixRule{
			baseRule{
				resourceType: "my_resource",
				attribute:    "my_attr",
			},
			tc.prefix,
		}

		actual := rule.Equates(tc.valueA, tc.valueB)
		if actual != tc.want {
			t.Errorf("Equates(%q, %q) = %t, want %t", tc.valueA, tc.valueB, actual, tc.want)
		}
	}
}
