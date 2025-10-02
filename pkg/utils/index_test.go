package utils

import "testing"

// TestTransformIndex tests various boundary and normal cases of the TransformIndex function.
// Test scenarios covered:
//   - Empty slice handling
//   - First element access
//   - In-range index conversion
//   - Lower bound clamping (index too small)
//   - Upper bound clamping (index too large)
func TestTransformIndex(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		indices := []int{-5, 0, 1, 5, 100}
		for _, idx := range indices {
			if got := TransformIndex([]int{}, idx); got != 0 {
				t.Errorf("empty slice should return 0 for index %d, got %d", idx, got)
			}
		}
	})

	type testCase struct {
		name  string   // test case name
		arr   []string // input array
		index int      // input index
		want  int      // expected result
	}

	cases := []testCase{
		{name: "first element", arr: []string{"a", "b", "c"}, index: 1, want: 0},
		{name: "in range", arr: []string{"a", "b", "c"}, index: 2, want: 1},
		{name: "clamp low", arr: []string{"a", "b", "c"}, index: 0, want: 0},
		{name: "clamp high", arr: []string{"a", "b", "c"}, index: 5, want: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := TransformIndex(tc.arr, tc.index); got != tc.want {
				t.Errorf("expected %d, got %d", tc.want, got)
			}
		})
	}
}
