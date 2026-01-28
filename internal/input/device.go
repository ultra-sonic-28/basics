package input

type Device interface {
	ReadLine() (string, error)
	GetChar() (rune, error)
}
