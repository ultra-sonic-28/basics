package common

type ForFrame[T any] struct {
	Var  string
	Data T
}

type ForStack[T any] struct {
	stack []ForFrame[T]
}

func NewForStack[T any]() *ForStack[T] {
	return &ForStack[T]{}
}

func (fs *ForStack[T]) Push(f ForFrame[T]) {
	fs.stack = append(fs.stack, f)
}

func (fs *ForStack[T]) Pop() (ForFrame[T], bool) {
	if len(fs.stack) == 0 {
		return ForFrame[T]{}, false
	}
	f := fs.stack[len(fs.stack)-1]
	fs.stack = fs.stack[:len(fs.stack)-1]
	return f, true
}

func (fs *ForStack[T]) IsEmpty() bool {
	return len(fs.stack) == 0
}
