package rules

import (
	"reflect"
	"testing"

	"github.com/busser/tfautomv/pkg/engine"
)

func TestParseRule(t *testing.T) {
	tests := []struct {
		s       string
		want    engine.Rule
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

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			actual, err := Parse(tt.s)

			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.wantErr {
				t.Errorf("expected error, got none")
			}
			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("ParseRule() mismatch:\ngot: %#v\nwant: %#v", actual, tt.want)
			}

			if s := actual.String(); s != tt.s {
				t.Errorf("String() = %q, want %q", s, tt.s)
			}
		})
	}
}
