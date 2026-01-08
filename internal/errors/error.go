package errors

import "fmt"

type Error struct {
	Kind   Kind
	Line   int
	Column int
	Token  string
	Msg    string
}

func (e *Error) Error() string {
	// Style Applesoft
	if e.Line > 0 {
		return fmt.Sprintf(
			"⚠️ %s IN %d (%s)",
			e.Msg,
			e.Line,
			e.Token,
		)
	}
	return fmt.Sprintf("⚠️ %s", e.Msg)
}

func NewSyntax(line, col int, token, msg string) *Error {
	return &Error{
		Kind:   Syntax,
		Line:   line,
		Column: col,
		Token:  token,
		Msg:    msg,
	}
}

func NewSemantic(line int, msg string) *Error {
	return &Error{
		Kind: Semantic,
		Line: line,
		Msg:  msg,
	}
}
