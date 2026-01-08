package interpreter

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

type Env struct {
	vars map[string]Value
}

func NewEnv() *Env {
	return &Env{
		vars: make(map[string]Value),
	}
}

func (e *Env) Set(name string, v Value) {
	e.vars[name] = v
}

func (e *Env) Get(name string) (Value, bool) {
	if v, ok := e.vars[name]; ok {
		return v, true
	}
	// Applesoft : variable non initialisée = 0
	return Value{Type: NUMBER, Num: 0}, false
}
