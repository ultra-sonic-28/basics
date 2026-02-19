package interpreter

import (
	"bytes"
	"testing"

	"basics/internal/common"
	"basics/internal/constants"
	"basics/internal/lexer"
	"basics/internal/machines"
	"basics/internal/parser"
	"basics/internal/runtime"
	"basics/testutils"
)

func TestInterpreter_DIM(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		assertFn func(t *testing.T, i *Interpreter)
		errorMsg string
	}{
		{
			name:   "DIM A(10)",
			source: `20 DIM A(10)`,
			assertFn: func(t *testing.T, i *Interpreter) {
				v, ok := i.rt.Env.Get("A")
				testutils.True(t, "A exists", ok)
				testutils.Equal(t, "A is ARRAY", v.Type, runtime.ARRAY)
				testutils.Equal(t, "base type", v.Array.BaseType, runtime.NUMBER)
				testutils.DeepEqual(t, "dimensions", v.Array.Dimensions, []int{11})
			},
		},
		{
			name:   "DIM A%(10)",
			source: `20 DIM A%(10)`,
			assertFn: func(t *testing.T, i *Interpreter) {
				v, _ := i.rt.Env.Get("A%")
				testutils.Equal(t, "base type", v.Array.BaseType, runtime.INTEGER)
				testutils.DeepEqual(t, "dimensions", v.Array.Dimensions, []int{11})
			},
		},
		{
			name:   "DIM A$(10)",
			source: `20 DIM A$(10)`,
			assertFn: func(t *testing.T, i *Interpreter) {
				v, _ := i.rt.Env.Get("A$")
				testutils.Equal(t, "base type", v.Array.BaseType, runtime.STRING)
				testutils.DeepEqual(t, "dimensions", v.Array.Dimensions, []int{11})
			},
		},
		{
			name: "multi-dim arrays",
			source: `
20 DIM A(10,20)
30 DIM A%(10, 20, 30)
40 DIM A$(1,2,3,4,5,6,7,8,9,10)
`,
			assertFn: func(t *testing.T, i *Interpreter) {
				v1, _ := i.rt.Env.Get("A")
				testutils.DeepEqual(t, "dims A", v1.Array.Dimensions, []int{11, 21})

				v2, _ := i.rt.Env.Get("A%")
				testutils.DeepEqual(t, "dims A%", v2.Array.Dimensions, []int{11, 21, 31})

				v3, _ := i.rt.Env.Get("A$")
				testutils.DeepEqual(t, "dims A$", v3.Array.Dimensions,
					[]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
				)
			},
		},
		{
			name:   "multiple arrays in one DIM",
			source: `20 DIM A(10), B%(50), C$(10, 10)`,
			assertFn: func(t *testing.T, i *Interpreter) {
				a, _ := i.rt.Env.Get("A")
				b, _ := i.rt.Env.Get("B%")
				c, _ := i.rt.Env.Get("C$")

				testutils.DeepEqual(t, "A dims", a.Array.Dimensions, []int{11})
				testutils.DeepEqual(t, "B dims", b.Array.Dimensions, []int{51})
				testutils.DeepEqual(t, "C dims", c.Array.Dimensions, []int{11, 11})
			},
		},
		{
			name: "DIM with expressions",
			source: `
20 Size = 50
30 DIM A%(2*Size, Size/10, 5)
`,
			assertFn: func(t *testing.T, i *Interpreter) {
				v, _ := i.rt.Env.Get("A%")
				testutils.DeepEqual(t, "dimensions",
					v.Array.Dimensions,
					[]int{101, 6, 6},
				)
			},
		},
		{
			name:     "error: dimension not numeric",
			source:   `20 DIM A$("HELLO")`,
			errorMsg: "⚠️ BAD SUBSCRIPT: DIMENSION MUST BE A NUMBER OR AN INTEGER\n",
		},
		{
			name:     "error: dimension zero",
			source:   `20 DIM A(0)`,
			errorMsg: "⚠️ BAD SUBSCRIPT: DIMENSION MUST BE POSITIVE\n",
		},
		{
			name:     "error: negative dimension",
			source:   `20 DIM A(-5)`,
			errorMsg: "⚠️ BAD SUBSCRIPT: DIMENSION MUST BE POSITIVE\n",
		},
	}

	rt, _ := machines.NewRuntime(constants.BASIC_TTY, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			rt.SetOutput(out)

			i := New(rt)

			tokens := lexer.Lex(tt.source)
			p := parser.New(tokens)
			prog, errs := p.ParseProgram()

			testutils.Equal(t, "no parser errors", len(errs), 0)

			testutils.Equal(t, "no parser errors", len(errs), 0)

			i.Run(prog)

			got := common.StripANSI(out.String())
			if tt.errorMsg != "" {
				testutils.Contains(t, "error message", got, tt.errorMsg)
				return
			}

			tt.assertFn(t, i)
		})
	}
}
