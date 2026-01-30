package runtime

import (
	"basics/testutils"
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
