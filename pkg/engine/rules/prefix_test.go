package rules

import "testing"

func TestPrefixRuleAppliesTo(t *testing.T) {
	rule := prefixRule{
		baseRule{
			resourceType: "my_resource",
			attribute:    "my_attr",
		},
		"does-not-matter",
	}

	tests := []struct {
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

	for _, tt := range tests {
		actual := rule.AppliesTo(tt.resourceType, tt.attribute)
		if actual != tt.want {
			t.Errorf("AppliesTo(%q, %q) = %t, want %t", tt.resourceType, tt.attribute, actual, tt.want)
		}
	}
}

func TestPrefixRuleEquates(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		rule := prefixRule{
			baseRule{
				resourceType: "my_resource",
				attribute:    "my_attr",
			},
			tt.prefix,
		}

		actual := rule.Equates(tt.valueA, tt.valueB)
		if actual != tt.want {
			t.Errorf("Equates(%q, %q) = %t, want %t", tt.valueA, tt.valueB, actual, tt.want)
		}
	}
}
