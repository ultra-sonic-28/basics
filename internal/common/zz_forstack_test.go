package common

import (
	"basics/testutils"
	"testing"
)

func TestForStack_PushPop_IsEmpty(t *testing.T) {
	// On teste un stack de int
	fs := NewForStack[int]()
	testutils.True(t, "stack initially empty", fs.IsEmpty())

	// Push un élément
	fs.Push(ForFrame[int]{Var: "i", Data: 42})
	testutils.False(t, "stack not empty after push", fs.IsEmpty())

	// Push un deuxième élément
	fs.Push(ForFrame[int]{Var: "j", Data: 99})
	testutils.False(t, "stack not empty after second push", fs.IsEmpty())

	// Pop le dernier élément
	f, ok := fs.Pop()
	testutils.True(t, "pop success", ok)
	testutils.Equal(t, "var name", f.Var, "j")
	testutils.Equal(t, "data value", f.Data, 99)
	testutils.False(t, "stack still not empty", fs.IsEmpty())

	// Pop le premier élément
	f, ok = fs.Pop()
	testutils.True(t, "pop success", ok)
	testutils.Equal(t, "var name", f.Var, "i")
	testutils.Equal(t, "data value", f.Data, 42)
	testutils.True(t, "stack empty after pops", fs.IsEmpty())

	// Pop sur stack vide
	_, ok = fs.Pop()
	testutils.False(t, "pop on empty stack returns false", ok)
}

func TestForStack_GenericString(t *testing.T) {
	// Stack de string
	fs := NewForStack[string]()
	fs.Push(ForFrame[string]{Var: "name", Data: "Alice"})
	f, ok := fs.Pop()
	testutils.True(t, "pop success", ok)
	testutils.Equal(t, "var name", f.Var, "name")
	testutils.Equal(t, "data value", f.Data, "Alice")
	testutils.True(t, "stack empty", fs.IsEmpty())
}
