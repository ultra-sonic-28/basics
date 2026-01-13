package interpreter

import (
	"fmt"
	"strings"

	"basics/internal/errors"
	"basics/internal/parser"
	"basics/internal/runtime"
)

//
// =======================
// Structures internes
// =======================
//

// Instruction représente un statement exécutable
type Instruction struct {
	LineNum int
	Stmt    parser.Statement
}

// ForFrame garde les infos d'une boucle FOR active
type ForFrame struct {
	Var     string
	End     float64
	Step    float64
	PCStart int // PC de l'instruction FOR
}

type ForStack struct {
	stack []ForFrame
}

func NewForStack() *ForStack {
	return &ForStack{}
}

func (fs *ForStack) Push(f ForFrame) {
	fs.stack = append(fs.stack, f)
}

func (fs *ForStack) Pop() {
	if len(fs.stack) > 0 {
		fs.stack = fs.stack[:len(fs.stack)-1]
	}
}

func (fs *ForStack) Top() *ForFrame {
	if len(fs.stack) == 0 {
		return nil
	}
	return &fs.stack[len(fs.stack)-1]
}

//
// =======================
// Interpreter
// =======================
//

type Interpreter struct {
	rt        *runtime.Runtime
	forStack  *ForStack
	insts     []Instruction
	lineIndex map[int]int // line number → PC
}

func New(rt *runtime.Runtime) *Interpreter {
	return &Interpreter{
		rt:       rt,
		forStack: NewForStack(),
	}
}

//
// =======================
// Programme → instructions
// =======================
//

func (i *Interpreter) buildInstructions(prog *parser.Program) {
	i.insts = nil
	i.lineIndex = make(map[int]int)

	for _, line := range prog.Lines {
		if _, exists := i.lineIndex[line.Number]; !exists {
			i.lineIndex[line.Number] = len(i.insts)
		}

		for _, stmt := range line.Stmts {
			i.insts = append(i.insts, Instruction{
				LineNum: line.Number,
				Stmt:    stmt,
			})
		}
	}
}

//
// =======================
// Boucle d'exécution
// =======================
//

func (i *Interpreter) Run(prog *parser.Program) {
	i.buildInstructions(prog)

	pc := 0
	for pc < len(i.insts) {
		inst := i.insts[pc]
		nextPC := pc + 1

		switch s := inst.Stmt.(type) {

		// -----------------------
		// END
		// -----------------------
		case *parser.EndStmt:
			i.rt.Halt()
			return

		// -----------------------
		// LET
		// -----------------------
		case *parser.LetStmt:
			val, err := EvalExpr(s.Value, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}
			i.rt.Env.Set(s.Name, val)

		// -----------------------
		// PRINT
		// -----------------------
		case *parser.PrintStmt:
			cursor := 0

			for iExpr, expr := range s.Exprs {
				val, err := EvalExpr(expr, i.rt)
				if err != nil {
					fmt.Println(err)
					return
				}

				str := ""
				switch val.Type {
				case runtime.NUMBER:
					str = formatNumber(val.Num)
				case runtime.STRING:
					str = val.Str
				}

				if iExpr > 0 {
					sep := s.Separators[iExpr-1]
					if sep == ',' {
						spaces := 14 - (cursor % 14)
						str = strings.Repeat(" ", spaces) + str
					}
				}

				i.rt.ExecPrint(str)
				cursor += len(str)
			}
			i.rt.ExecPrint("\n")

		// -----------------------
		// HTAB / VTAB
		// -----------------------
		case *parser.HTabStmt:
			val, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}
			i.rt.ExecHTab(int(val.Num))

		case *parser.VTabStmt:
			val, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}
			i.rt.ExecVTab(int(val.Num))

		// -----------------------
		// FOR
		// -----------------------
		case *parser.ForStmt:
			startVal, err := EvalExpr(s.Start, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}
			endVal, err := EvalExpr(s.End, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}

			step := 1.0
			if s.Step != nil {
				stepVal, err := EvalExpr(s.Step, i.rt)
				if err != nil {
					fmt.Println(err)
					return
				}
				step = stepVal.Num
				if step == 0 {
					fmt.Println(errors.NewSemantic(
						inst.LineNum,
						"STEP CANNOT BE ZERO",
					))
					return
				}
			}

			i.rt.Env.Set(s.Var, runtime.Value{
				Type: runtime.NUMBER,
				Num:  startVal.Num,
			})

			i.forStack.Push(ForFrame{
				Var:     s.Var,
				End:     endVal.Num,
				Step:    step,
				PCStart: pc,
			})

		// -----------------------
		// NEXT
		// -----------------------
		case *parser.NextStmt:
			frame := i.forStack.Top()
			if frame == nil {
				fmt.Println("?NEXT WITHOUT FOR")
				return
			}

			v, _ := i.rt.Env.Get(frame.Var)
			v.Num += frame.Step

			done := (frame.Step > 0 && v.Num > frame.End) ||
				(frame.Step < 0 && v.Num < frame.End)

			if !done {
				i.rt.Env.Set(frame.Var, v)
				nextPC = frame.PCStart + 1
			} else {
				i.forStack.Pop()
			}

		// -----------------------
		// GOTO
		// -----------------------
		case *parser.GotoStmt:
			val, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}

			if val.Type != runtime.NUMBER {
				fmt.Println("?GOTO TYPE MISMATCH")
				return
			}

			line := int(val.Num)
			targetPC, ok := i.lineIndex[line]
			if !ok {
				fmt.Printf("?UNDEFINED LINE %d\n", line)
				return
			}

			nextPC = targetPC
		}

		pc = nextPC
	}
}

//
// =======================
// Utils
// =======================
//

func formatNumber(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%g", f)
}
