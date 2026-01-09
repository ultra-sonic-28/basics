package runtime

import "fmt"

type ValueType int

const (
	NUMBER ValueType = iota
	STRING
)

type Value struct {
	Type ValueType
	Num  float64
	Str  string
}

type Environment struct {
	vars map[string]Value
}

func NewEnvironment() *Environment {
	return &Environment{
		vars: make(map[string]Value),
	}
}

func (e *Environment) Set(name string, v Value) {
	e.vars[name] = v
}

func (e *Environment) Get(name string) (Value, bool) {
	if v, ok := e.vars[name]; ok {
		return v, true
	}
	// Applesoft : variable non initialisée = 0
	return Value{Type: NUMBER, Num: 0}, false
}

func (v Value) String() string {
	if v.Type == STRING {
		return v.Str
	}
	return fmt.Sprintf("%f", v.Num)
}
