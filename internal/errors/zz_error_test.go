package errors

import (
	"testing"

	"basics/testutils"
)

func TestError_Error_WithLine(t *testing.T) {
	err := &Error{
		Kind:   Syntax,
		Line:   120,
		Column: 5,
		Token:  "PRINT",
		Msg:    "SYNTAX ERROR",
	}

	got := err.Error()
	want := "⚠️ SYNTAX ERROR IN 120 (PRINT)"

	testutils.Equal(t, "", got, want)
}

func TestError_Error_WithoutLine(t *testing.T) {
	err := &Error{
		Kind: Semantic,
		Msg:  "DIVISION BY ZERO",
	}

	got := err.Error()
	want := "⚠️ DIVISION BY ZERO"

	testutils.Equal(t, "", got, want)
}

func TestNewSyntax(t *testing.T) {
	err := NewSyntax(10, 3, "IF", "SYNTAX ERROR")

	testutils.Equal(t, "", err.Kind, Syntax)
	testutils.Equal(t, "", err.Line, 10)
	testutils.Equal(t, "", err.Column, 3)
	testutils.Equal(t, "", err.Token, "IF")
	testutils.Equal(t, "", err.Msg, "SYNTAX ERROR")
}

func TestNewSemantic(t *testing.T) {
	err := NewSemantic(200, "TYPE MISMATCH")

	testutils.Equal(t, "", err.Kind, Semantic)
	testutils.Equal(t, "", err.Line, 200)
	testutils.Equal(t, "", err.Msg, "TYPE MISMATCH")

	// Champs non utilisés
	testutils.Equal(t, "", err.Column, 0)
	testutils.Equal(t, "", err.Token, "")
}
