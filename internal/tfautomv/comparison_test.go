package tfautomv

import (
	"sort"
	"testing"

	"github.com/padok-team/tfautomv/internal/slices"
	"github.com/padok-team/tfautomv/internal/tfautomv/ignore"
)

func TestCompare(t *testing.T) {
	tt := []struct {
		name            string
		created         *Resource
		destroyed       *Resource
		rules           []ignore.Rule
		wantMatching    []string
		wantIgnored     []string
		wantMismatching []string
	}{
		{
			name: "without rules",
			created: &Resource{
				Attributes: map[string]interface{}{
					"a": "hello",
					"b": 123,
					"c": true,
					"d": nil,
					"e": "foo",
					"f": 456,
					"h": "goodbye",
				},
			},
			destroyed: &Resource{
				Attributes: map[string]interface{}{
					"a": "hello",
					"b": 123,
					"c": false,
					"d": 12.34,
					"e": nil,
					"g": "whatever",
					"h": 789,
				},
			},
			wantMatching:    []string{"a", "b"},
			wantMismatching: []string{"c", "e", "f", "h"},
		},
		{
			name: "with rules",
			created: &Resource{
				Type: "my_resource",
				Attributes: map[string]interface{}{
					"a": "hello",
					"b": 123,
					"c": true,
					"d": nil,
					"e": "foo",
					"f": 456,
					"h": "goodbye",
					"i": "{\"foo\":\"bar\"}",
				},
			},
			destroyed: &Resource{
				Type: "my_resource",
				Attributes: map[string]interface{}{
					"a": "hello",
					"b": 123,
					"c": false,
					"d": 12.34,
					"e": nil,
					"g": "whatever",
					"h": 789,
					"i": "{\n\t\"foo\": \"bar\"\n}",
				},
			},
			rules: []ignore.Rule{
				ignore.MustParseRule("everything:my_resource:c"),
				ignore.MustParseRule("whitespace:my_resource:i"),
			},
			wantMatching:    []string{"a", "b"},
			wantIgnored:     []string{"c", "i"},
			wantMismatching: []string{"e", "f", "h"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := Compare(tc.created, tc.destroyed, tc.rules)
			if actual.Created != tc.created {
				t.Errorf("Compare(): resulting Comparison does not point to created resource")
			}
			if actual.Destroyed != tc.destroyed {
				t.Errorf("Compare(): resulting Comparison does not point to destroyed resource")
			}

			sort.Strings(actual.MatchingAttributes)
			sort.Strings(tc.wantMatching)
			if !slices.Equal(actual.MatchingAttributes, tc.wantMatching) {
				t.Errorf("Compare().MatchingAttributes = %#v, want %#v",
					actual.MatchingAttributes, tc.wantMatching)
			}

			sort.Strings(actual.IgnoredAttributes)
			sort.Strings(tc.wantIgnored)
			if !slices.Equal(actual.IgnoredAttributes, tc.wantIgnored) {
				t.Errorf("Compare().IgnoredAttributes = %#v, want %#v",
					actual.IgnoredAttributes, tc.wantIgnored)
			}

			sort.Strings(actual.MismatchingAttributes)
			sort.Strings(tc.wantMismatching)
			if !slices.Equal(actual.MismatchingAttributes, tc.wantMismatching) {
				t.Errorf("Compare().MismatchingAttributes = %#v, want %#v",
					actual.MismatchingAttributes, tc.wantMismatching)
			}
		})
	}
}

func TestIsMatch(t *testing.T) {
	tt := []struct {
		comp Comparison
		want bool
	}{
		{
			comp: Comparison{
				MatchingAttributes:    []string{"a", "b", "c"},
				MismatchingAttributes: nil,
			},
			want: true,
		},
		{
			comp: Comparison{
				MatchingAttributes:    []string{"a", "b", "c"},
				MismatchingAttributes: []string{"d"},
			},
			want: false,
		},
	}

	for _, tc := range tt {
		actual := tc.comp.IsMatch()
		if actual != tc.want {
			t.Errorf("IsMatch() = %#v, want %#v when MatchingAttributes = %#v and MismatchingAttributes = %#v",
				actual, tc.want, tc.comp.MatchingAttributes, tc.comp.MismatchingAttributes)
		}
	}
}
