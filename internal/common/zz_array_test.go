package common

import (
	"basics/testutils"
	"fmt"
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []int
	}{
		{
			name:     "single int",
			input:    5,
			expected: []int{5},
		},
		{
			name:     "simple slice",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "nested slices",
			input:    []interface{}{1, []int{2, 3}, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			name: "deeply nested mixed slices",
			input: []interface{}{
				1,
				[]interface{}{
					2,
					[]int{3, 4},
				},
				5,
			},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "array input",
			input:    [3]int{7, 8, 9},
			expected: []int{7, 8, 9},
		},
		{
			name:     "mixed array and slice",
			input:    []interface{}{[2]int{1, 2}, []int{3, 4}},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "unsupported types ignored",
			input:    []interface{}{1, "foo", 2.5, []int{2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.input)

			str := fmt.Sprintf(
				"Flatten(%#v): got %v, expected %v",
				tt.input,
				got,
				tt.expected,
			)
			testutils.True(t, str, reflect.DeepEqual(got, tt.expected))
		})
	}
}
