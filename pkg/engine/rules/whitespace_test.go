package rules

import "testing"

func TestWhitespaceRuleAppliesTo(t *testing.T) {
	rule := whitespaceRule{
		baseRule{
			resourceType: "my_resource",
			attribute:    "my_attr",
		},
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

func TestWhitespaceRuleEquates(t *testing.T) {
	rule := whitespaceRule{
		baseRule{
			resourceType: "my_resource",
			attribute:    "my_attr",
		},
	}

	tests := []struct {
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

	for _, tt := range tests {
		actual := rule.Equates(tt.valueA, tt.valueB)
		if actual != tt.want {
			t.Errorf("Equates(%q, %q) = %t, want %t", tt.valueA, tt.valueB, actual, tt.want)
		}
	}
}

func TestWithoutWhitespace(t *testing.T) {
	tests := []struct {
		str  string
		want string
	}{
		{"", ""},
		{"foobar", "foobar"},
		{"  foo  ", "foo"},
		{"\t foo\tbar\n", "foobar"},
	}

	for _, tt := range tests {
		actual := withoutWhitespace(tt.str)
		if actual != tt.want {
			t.Errorf("withoutWhitespace(%q) = %q, want %q", tt.str, actual, tt.want)
		}
	}
}
