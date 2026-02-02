package runtime

import (
	"basics/testutils"
	"fmt"
	"testing"
)

func TestEnvironment_Clear_ResetValues(t *testing.T) {
	env := NewEnvironment()

	env.Set("A", Value{Type: NUMBER, Num: 42})
	env.Set("I%", Value{Type: INTEGER, Int: 7})
	env.Set("S$", Value{Type: STRING, Str: "HELLO"})
	env.Set("B", Value{Type: BOOLEAN, Flag: true})

	env.Clear()

	v, _ := env.Get("A")
	testutils.Equal(t, "A reset", v.Num, 0.0)

	v, _ = env.Get("I%")
	testutils.Equal(t, "I% reset", v.Int, 0)

	v, _ = env.Get("S$")
	testutils.Equal(t, "S$ reset", v.Str, "")

	v, _ = env.Get("B")
	testutils.Equal(t, "B reset", v.Flag, false)
}

func TestEnvironment_Clear_ResetArrays(t *testing.T) {
	tests := []struct {
		name       string
		arrayName  string
		baseType   ValueType
		dims       []int
		fillFn     func(a *ArrayValue)
		assertZero func(t *testing.T, a *ArrayValue)
	}{
		{
			name:      "integer array 1D",
			arrayName: "A%",
			baseType:  INTEGER,
			dims:      []int{5},
			fillFn: func(a *ArrayValue) {
				for i := range a.Data {
					a.Data[i].Int = i + 1
				}
			},
			assertZero: func(t *testing.T, a *ArrayValue) {
				for _, v := range a.Data {
					testutils.Equal(t, "cell reset", v.Int, 0)
					testutils.Equal(t, "cell type", v.Type, INTEGER)
				}
			},
		},
		{
			name:      "number array 2D",
			arrayName: "B",
			baseType:  NUMBER,
			dims:      []int{3, 4},
			fillFn: func(a *ArrayValue) {
				for i := range a.Data {
					a.Data[i].Num = float64(i) + 0.5
				}
			},
			assertZero: func(t *testing.T, a *ArrayValue) {
				for _, v := range a.Data {
					testutils.Equal(t, "cell reset", v.Num, 0.0)
					testutils.Equal(t, "cell type", v.Type, NUMBER)
				}
			},
		},
		{
			name:      "string array 1D",
			arrayName: "S$",
			baseType:  STRING,
			dims:      []int{10},
			fillFn: func(a *ArrayValue) {
				for i := range a.Data {
					a.Data[i].Str = "TEST"
				}
			},
			assertZero: func(t *testing.T, a *ArrayValue) {
				for _, v := range a.Data {
					testutils.Equal(t, "cell reset", v.Str, "")
					testutils.Equal(t, "cell type", v.Type, STRING)
				}
			},
		},
		{
			name:      "integer array 3D",
			arrayName: "C%",
			baseType:  INTEGER,
			dims:      []int{2, 2, 2},
			fillFn: func(a *ArrayValue) {
				for i := range a.Data {
					a.Data[i].Int = 99
				}
			},
			assertZero: func(t *testing.T, a *ArrayValue) {
				for _, v := range a.Data {
					testutils.Equal(t, "cell reset", v.Int, 0)
					testutils.Equal(t, "cell type", v.Type, INTEGER)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()

			arr := NewArray(tt.baseType, tt.dims)
			tt.fillFn(arr)

			env.Set(tt.arrayName, Value{
				Type:  ARRAY,
				Array: arr,
			})

			env.Clear()

			v, ok := env.Get(tt.arrayName)
			testutils.True(t, fmt.Sprintf("test[%s] - array still exists", tt.name), ok)
			testutils.Equal(t, fmt.Sprintf("test[%s] - value type is ARRAY", tt.name), v.Type, ARRAY)

			a := v.Array
			testutils.Equal(t, fmt.Sprintf("test[%s] - base type preserved", tt.name), a.BaseType, tt.baseType)
			testutils.DeepEqual(t, fmt.Sprintf("test[%s] - dimensions preserved", tt.name), a.Dimensions, tt.dims)

			tt.assertZero(t, a)
		})
	}
}
