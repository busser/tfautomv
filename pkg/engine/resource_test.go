package engine_test

import (
	"testing"

	"github.com/busser/tfautomv/pkg/engine"
	"github.com/busser/tfautomv/pkg/engine/rules"

	"github.com/google/go-cmp/cmp"
)

func TestCompareResources(t *testing.T) {
	tests := []struct {
		name string

		create engine.Resource
		delete engine.Resource
		rules  []engine.Rule

		wantMatching    []string
		wantMismatching []string
		wantIgnored     []string
	}{
		{
			name: "without rules",
			create: dummyResource(map[string]any{
				"a": "hello",
				"b": 123,
				"c": true,
				"d": nil,
				"e": "foo",
				"f": 456,
				"h": "goodbye",
			}),
			delete: dummyResource(map[string]any{
				"a": "hello",
				"b": 123,
				"c": false,
				"d": 12.34,
				"e": nil,
				"g": "whatever",
				"h": 789,
			}),
			wantMatching:    []string{"a", "b"},
			wantMismatching: []string{"c", "e", "f", "h"},
		},

		{
			name: "with rules",
			create: dummyResource(map[string]any{
				"a": "hello",
				"b": 123,
				"c": true,
				"d": nil,
				"e": "foo",
				"f": 456,
				"h": "goodbye",
				"i": "{\"foo\":\"bar\"}",
				"j": "some_string",
			}),
			delete: dummyResource(map[string]any{
				"a": "hello",
				"b": 123,
				"c": false,
				"d": 12.34,
				"e": nil,
				"g": "whatever",
				"h": 789,
				"i": "{\n\t\"foo\": \"bar\"\n}",
				"j": "b/some_string",
			}),
			rules: []engine.Rule{
				rules.MustParse("everything:dummy_type:c"),
				rules.MustParse("whitespace:dummy_type:i"),
				rules.MustParse("prefix:dummy_type:j:b/"),
			},
			wantMatching:    []string{"a", "b"},
			wantMismatching: []string{"e", "f", "h"},
			wantIgnored:     []string{"c", "i", "j"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := engine.ResourceComparison{
				ToCreate:              tt.create,
				ToDelete:              tt.delete,
				MatchingAttributes:    tt.wantMatching,
				MismatchingAttributes: tt.wantMismatching,
				IgnoredAttributes:     tt.wantIgnored,
			}
			got := engine.CompareResources(tt.create, tt.delete, tt.rules)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

func dummyResource(attributes map[string]any) engine.Resource {
	return engine.Resource{
		ModuleID:   "dummy_module_id",
		Type:       "dummy_type",
		Address:    "dummy_address",
		Attributes: attributes,
	}
}
