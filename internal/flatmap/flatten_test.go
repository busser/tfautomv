package flatmap_test

import (
	"testing"

	"github.com/busser/tfautomv/internal/flatmap"

	"github.com/google/go-cmp/cmp"
)

func TestToString(t *testing.T) {
	tt := []struct {
		name string
		obj  map[string]interface{}
		want map[string]interface{}
	}{
		{
			name: "already-flat",
			obj: map[string]interface{}{
				"string": "bar",
				"int":    123,
				"float":  1.23,
				"bool":   true,
				"nil":    nil,
			},
			want: map[string]interface{}{
				"string": "bar",
				"int":    123,
				"float":  1.23,
				"bool":   true,
				"nil":    nil,
			},
		},
		{
			name: "slice",
			obj: map[string]interface{}{
				"ints": []int{1, 1, 2, 3, 5, 8, 13},
			},
			want: map[string]interface{}{
				"ints.#": 7,
				"ints.0": 1,
				"ints.1": 1,
				"ints.2": 2,
				"ints.3": 3,
				"ints.4": 5,
				"ints.5": 8,
				"ints.6": 13,
			},
		},
		{
			name: "map",
			obj: map[string]interface{}{
				"map": map[string]int{
					"foo": 123,
					"bar": 456,
				},
			},
			want: map[string]interface{}{
				"map.foo": 123,
				"map.bar": 456,
			},
		},
		{
			name: "complex",
			obj: map[string]interface{}{
				"string": "bar",
				"map": map[string]interface{}{
					"int":   0,
					"float": 1.23,
					"map": map[string]interface{}{
						"string": "foo",
						"int":    123,
						"slice":  []int{1, 1, 2, 3, 5, 8, 13},
					},
					"slice": []interface{}{
						false,
						map[string]int{
							"foo": 123,
						},
					},
					"nil": nil,
				},
				"bool": true,
			},
			want: map[string]interface{}{
				"string":          "bar",
				"map.int":         0,
				"map.float":       1.23,
				"map.map.string":  "foo",
				"map.map.int":     123,
				"map.map.slice.#": 7,
				"map.map.slice.0": 1,
				"map.map.slice.1": 1,
				"map.map.slice.2": 2,
				"map.map.slice.3": 3,
				"map.map.slice.4": 5,
				"map.map.slice.5": 8,
				"map.map.slice.6": 13,
				"map.slice.#":     2,
				"map.slice.0":     false,
				"map.slice.1.foo": 123,
				"map.nil":         nil,
				"bool":            true,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := flatmap.Flatten(tc.obj)
			if err != nil {
				t.Fatalf("ToString() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ToString() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
