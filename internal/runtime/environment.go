package runtime

import "fmt"

type ValueType int

const (
	NUMBER ValueType = iota
	INTEGER
	STRING
	BOOLEAN
	ARRAY
)

type ArrayValue struct {
	BaseType   ValueType // NUMBER / INTEGER / STRING
	Dimensions []int     // tailles déclarées (+1 déjà appliqué)
	Data       []Value   // stockage linéaire
}

type Value struct {
	Type  ValueType
	Num   float64
	Int   int
	Str   string
	Flag  bool
	Array *ArrayValue
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

	if v.Type == INTEGER {
		return fmt.Sprintf("%d", v.Int)
	}

	return fmt.Sprintf("%f", v.Num)
}

func (e *Environment) Clear() {
	for name, v := range e.vars {
		switch v.Type {
		case NUMBER:
			v.Num = 0

		case INTEGER:
			v.Int = 0

		case STRING:
			v.Str = ""

		case BOOLEAN:
			v.Flag = false
		}

		e.vars[name] = v
	}
}

func (a *ArrayValue) index(indices []int) (int, error) {
	if len(indices) != len(a.Dimensions) {
		return 0, fmt.Errorf("BAD SUBSCRIPT")
	}

	stride := 1
	idx := 0

	for i := len(a.Dimensions) - 1; i >= 0; i-- {
		if indices[i] < 0 || indices[i] >= a.Dimensions[i] {
			return 0, fmt.Errorf("BAD SUBSCRIPT")
		}
		idx += indices[i] * stride
		stride *= a.Dimensions[i]
	}

	return idx, nil
}

func NewArray(baseType ValueType, dims []int) *ArrayValue {
	total := 1
	for _, d := range dims {
		total *= d
	}

	data := make([]Value, total)

	for i := range data {
		switch baseType {
		case INTEGER:
			data[i] = Value{Type: INTEGER, Int: 0}
		case STRING:
			data[i] = Value{Type: STRING, Str: ""}
		default:
			data[i] = Value{Type: NUMBER, Num: 0}
		}
	}

	return &ArrayValue{
		BaseType:   baseType,
		Dimensions: dims,
		Data:       data,
	}
}
