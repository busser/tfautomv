package ignore

import "testing"

func TestWhitespaceRuleAppliesTo(t *testing.T) {
	rule := whitespaceRule{
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

func TestWhitespaceRuleEquates(t *testing.T) {
	rule := whitespaceRule{
		baseRule{
			resourceType: "my_resource",
			attribute:    "my_attr",
		},
	}

	tt := []struct {
		valueA interface{}
		valueB interface{}
		want   bool
	}{
		{
			valueA: "foo",
			valueB: "foo",
			want:   true,
		},
		{
			valueA: "\tfoo bar\n",
			valueB: "foo\tbar",
			want:   true,
		},
		{
			valueA: 123,
			valueB: 123,
			want:   true,
		},
		{
			valueA: 123,
			valueB: 456,
			want:   false,
		},
		{
			valueA: 123,
			valueB: "123",
			want:   false,
		},
		{
			valueA: false,
			valueB: false,
			want:   true,
		},
		{
			valueA: false,
			valueB: "false",
			want:   false,
		},
	}

	for _, tc := range tt {
		actual := rule.Equates(tc.valueA, tc.valueB)
		if actual != tc.want {
			t.Errorf("Equates(%q, %q) = %t, want %t", tc.valueA, tc.valueB, actual, tc.want)
		}
	}
}

func TestWithoutWhitespace(t *testing.T) {
	tt := []struct {
		str  string
		want string
	}{
		{"", ""},
		{"foobar", "foobar"},
		{"  foo  ", "foo"},
		{"\t foo\tbar\n", "foobar"},
	}

	for _, tc := range tt {
		actual := withoutWhitespace(tc.str)
		if actual != tc.want {
			t.Errorf("withoutWhitespace(%q) = %q, want %q", tc.str, actual, tc.want)
		}
	}
}
