package input

import (
	"basics/internal/logger"
)

type FakeInput struct {
	str string
}

func NewFakeInput(s string) *FakeInput {
	logger.Info("Instanciate FakeTTY renderer")
	return &FakeInput{
		str: s,
	}
}

func (f *FakeInput) GetChar() (rune, error) {
	var r rune
	for _, ch := range f.str { // le premier passage donne la première rune
		r = ch
		break
	}
	return r, nil

}

func (f *FakeInput) ReadLine() (string, error) {
	return f.str, nil
}
