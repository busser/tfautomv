package ignore

import (
	"reflect"
	"testing"
)

func TestParseRule(t *testing.T) {
	tt := []struct {
		s       string
		want    Rule
		wantErr bool
	}{
		// Everything rule
		{
			s: "everything:my_resource:my_attr",
			want: &everythingRule{
				baseRule{
					resourceType: "my_resource",
					attribute:    "my_attr",
				},
			},
		},
		{
			s:       "everything:my_resource",
			wantErr: true,
		},
		{
			s:       "everything:my_resource:my_attr:extra",
			wantErr: true,
		},

		// Whitespace rule
		{
			s: "whitespace:my_resource:my_attr",
			want: &whitespaceRule{
				baseRule{
					resourceType: "my_resource",
					attribute:    "my_attr",
				},
			},
		},
		{
			s:       "whitespace:my_resource",
			wantErr: true,
		},
		{
			s:       "whitespace:my_resource:my_attr:extra",
			wantErr: true,
		},

		// Non-existent rule
		{
			s:       "doesnotexist:foo:bar",
			wantErr: true,
		},
		{
			s:       "tooshort",
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.s, func(t *testing.T) {
			actual, err := ParseRule(tc.s)

			if err != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tc.wantErr {
				t.Errorf("expected error, got none")
			}
			if tc.wantErr {
				return
			}

			if !reflect.DeepEqual(actual, tc.want) {
				t.Errorf("ParseRule() mismatch:\ngot: %#v\nwant: %#v", actual, tc.want)
			}

			if s := actual.String(); s != tc.s {
				t.Errorf("String() = %q, want %q", s, tc.s)
			}
		})
	}
}
