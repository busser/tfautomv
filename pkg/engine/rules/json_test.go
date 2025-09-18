package rules

import "testing"

func TestJSONRuleAppliesTo(t *testing.T) {
	rule := jsonRule{
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

func TestJSONRuleEquates(t *testing.T) {
	tests := []struct {
		valueA interface{}
		valueB interface{}
		prefix string
		want   bool
	}{
		{
			valueA: "foo",
			valueB: "foo",
			want:   false,
		},
		{
			valueA: `"foo"`,
			valueB: `"foo"`,
			want:   true,
		},
		{
			valueA: "[1,2,3]",
			valueB: "[1,2,3]",
			want:   true,
		},
		{
			valueA: "[1,2,3]",
			valueB: "[1,3,2]",
			want:   false,
		},
		{
			valueA: `{"foo": "bar", "baz": "qux"}`,
			valueB: `{"baz": "qux", "foo": "bar"}`,
			want:   true,
		},
		{
			valueA: `{"foo": {"list": [1,2,3], "object": {"nested": true}}}`,
			valueB: `{
        "foo": {
          "object": {
            "nested": true
          },
         "list": [
            1,
            2,
            3
          ]
        }
      }`,
			want: true,
		},
	}

	for _, tt := range tests {
		rule := jsonRule{
			baseRule{
				resourceType: "my_resource",
				attribute:    "my_attr",
			},
		}

		actual := rule.Equates(tt.valueA, tt.valueB)
		if actual != tt.want {
			t.Errorf("Equates(%q, %q) = %t, want %t", tt.valueA, tt.valueB, actual, tt.want)
		}
	}
}
