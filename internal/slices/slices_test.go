package slices

import (
	"testing"
)

func TestIndex(t *testing.T) {
	tt := []struct {
		s    []int
		v    int
		want int
	}{
		{[]int{1, 1, 2, 3, 5}, 0, -1},
		{[]int{1, 1, 2, 3, 5}, 1, 0},
		{[]int{1, 1, 2, 3, 5}, 2, 2},
		{[]int{1, 1, 2, 3, 5}, 3, 3},
		{[]int{1, 1, 2, 3, 5}, 4, -1},
		{[]int{1, 1, 2, 3, 5}, 5, 4},
		{[]int{}, 123, -1},
	}

	for _, tc := range tt {
		actual := Index(tc.s, tc.v)
		if actual != tc.want {
			t.Errorf("Index(%#v, %#v) = %#v, want %#v",
				tc.s, tc.s, actual, tc.want)
		}
	}
}

func TestContains(t *testing.T) {
	tt := []struct {
		s    []int
		v    int
		want bool
	}{
		{[]int{1, 1, 2, 3, 5}, 0, false},
		{[]int{1, 1, 2, 3, 5}, 1, true},
		{[]int{1, 1, 2, 3, 5}, 2, true},
		{[]int{1, 1, 2, 3, 5}, 3, true},
		{[]int{1, 1, 2, 3, 5}, 4, false},
		{[]int{1, 1, 2, 3, 5}, 5, true},
		{[]int{}, 123, false},
	}

	for _, tc := range tt {
		actual := Contains(tc.s, tc.v)
		if actual != tc.want {
			t.Errorf("Contains(%#v, %#v) = %#v, want %#v",
				tc.s, tc.s, actual, tc.want)
		}
	}
}

func TestEqual(t *testing.T) {
	tt := []struct {
		a, b []int
		want bool
	}{
		{[]int{1, 2, 3}, []int{1, 2, 3}, true},
		{[]int{1, 2, 3}, []int{1, 2, 4}, false},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3}, false},
		{[]int{}, []int{1, 2, 3}, false},
		{nil, []int{1, 2, 3}, false},
		{nil, nil, true},
		{nil, []int{}, true},
		{[]int{}, []int{}, true},
	}

	for _, tc := range tt {
		actual := Equal(tc.a, tc.b)
		if actual != tc.want {
			t.Errorf("Equal(%#v, %#v) = %#v, want %#v",
				tc.a, tc.b, actual, tc.want)
		}

		// Equality is symetric.
		actual = Equal(tc.b, tc.a)
		if actual != tc.want {
			t.Errorf("Equal(%#v, %#v) = %#v, want %#v",
				tc.b, tc.a, actual, tc.want)
		}
	}
}
