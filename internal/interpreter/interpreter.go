package interpreter

import (
	"fmt"
	"strings"

	"basics/internal/errors"
	"basics/internal/logger"
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

// For GOSUB ... RETURN
type GosubFrame struct {
	ReturnPC int
}

type GosubStack struct {
	stack []GosubFrame
}

func NewGosubStack() *GosubStack {
	return &GosubStack{}
}

func (s *GosubStack) Push(pc int) {
	s.stack = append(s.stack, GosubFrame{ReturnPC: pc})
}

func (s *GosubStack) Pop() (int, bool) {
	if len(s.stack) == 0 {
		return 0, false
	}
	top := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return top.ReturnPC, true
}

//
// =======================
// Interpreter
// =======================
//

type Interpreter struct {
	rt         *runtime.Runtime
	forStack   *ForStack
	gosubStack *GosubStack
	insts      []Instruction
	lineIndex  map[int]int // line number → PC
}

func New(rt *runtime.Runtime) *Interpreter {
	return &Interpreter{
		rt:         rt,
		forStack:   NewForStack(),
		gosubStack: NewGosubStack(),
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
	logger.Debug(fmt.Sprintf("Program contains %d instructions", len(i.insts)))

	pc := 0
	for pc < len(i.insts) {
		inst := i.insts[pc]
		nextPC := pc + 1

		logger.Debug(fmt.Sprintf("Executing line: %d, instruction: %d", inst.LineNum, pc))

		switch s := inst.Stmt.(type) {

		// -----------------------
		// HOME
		// -----------------------
		case *parser.HomeStmt:
			i.rt.ExecHome()

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
		// FOR (Applesoft semantics)
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

			end := float64(int(endVal.Num + 0.5))

			// 🔹 Initialisation TOUJOURS faite
			i.rt.Env.Set(s.Var, runtime.Value{
				Type: runtime.NUMBER,
				Num:  startVal.Num,
			})

			// 🔹 Empiler SANS TEST
			i.forStack.Push(ForFrame{
				Var:     s.Var,
				End:     end,
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

		// -----------------------
		// GOSUB
		// -----------------------
		case *parser.GosubStmt:
			val, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}

			if val.Type != runtime.NUMBER {
				fmt.Println("?GOSUB TYPE MISMATCH")
				return
			}

			line := int(val.Num)
			targetPC, ok := i.lineIndex[line]
			if !ok {
				fmt.Printf("?UNDEFINED LINE %d\n", line)
				return
			}

			// ⚠️ empiler l’instruction SUIVANTE
			i.gosubStack.Push(pc + 1)

			nextPC = targetPC

		// -----------------------
		// RETURN
		// -----------------------
		case *parser.ReturnStmt:
			retPC, ok := i.gosubStack.Pop()
			if !ok {
				fmt.Println("?RETURN WITHOUT GOSUB")
				return
			}
			nextPC = retPC

		// -----------------------
		// IF
		// -----------------------
		case *parser.IfStmt:
			cond, err := EvalExpr(s.Cond, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}

			exec := false

			switch cond.Type {
			case runtime.BOOLEAN:
				exec = cond.Flag
			case runtime.NUMBER:
				exec = cond.Num != 0
			}

			if exec {
				for _, stmt := range s.Then {
					nextPC = i.execInline(inst.LineNum, stmt, pc)
				}
			} else if s.Else != nil {
				for _, stmt := range s.Else {
					nextPC = i.execInline(inst.LineNum, stmt, pc)
				}
			}
		}

		pc = nextPC
	}
}

// =======================
// Inline execution helper
// =======================

func (i *Interpreter) execInline(line int, stmt parser.Statement, pc int) int {
	switch s := stmt.(type) {

	case *parser.HomeStmt:
		i.rt.ExecHome()
		return pc + 1

	case *parser.GotoStmt:
		val, err := EvalExpr(s.Expr, i.rt)
		if err != nil {
			fmt.Println(err)
			return pc + 1
		}
		target, ok := i.lineIndex[int(val.Num)]
		if !ok {
			fmt.Printf("?UNDEFINED LINE %d\n", int(val.Num))
			return pc + 1
		}
		return target

	case *parser.GosubStmt:
		val, err := EvalExpr(s.Expr, i.rt)
		if err != nil {
			fmt.Println(err)
			return pc + 1
		}

		line := int(val.Num)
		targetPC, ok := i.lineIndex[line]
		if !ok {
			fmt.Printf("?UNDEFINED LINE %d\n", line)
			return pc + 1
		}

		i.gosubStack.Push(pc + 1)
		return targetPC

	case *parser.ReturnStmt:
		retPC, ok := i.gosubStack.Pop()
		if !ok {
			fmt.Println("?RETURN WITHOUT GOSUB")
			return pc + 1
		}
		return retPC

	case *parser.LetStmt:
		val, err := EvalExpr(s.Value, i.rt)
		if err != nil {
			fmt.Println(err)
			return pc + 1
		}
		i.rt.Env.Set(s.Name, val)
		return pc + 1

	case *parser.PrintStmt:
		//i.rt.ExecPrint(s.Exprs[0].(*parser.StringLiteral).Value)
		//i.rt.ExecPrint("\n")
		for iExpr, expr := range s.Exprs {
			val, err := EvalExpr(expr, i.rt)
			if err != nil {
				fmt.Println(err)
				return pc
			}

			var out string
			switch val.Type {
			case runtime.STRING:
				out = val.Str
			case runtime.NUMBER:
				out = formatNumber(val.Num)
			case runtime.BOOLEAN:
				if val.Flag {
					out = "1"
				} else {
					out = "0"
				}
			default:
				out = ""
			}

			// séparateurs ; et ,
			if iExpr > 0 {
				sep := s.Separators[iExpr-1]
				switch sep {
				case ',':
					out = " " + out
				case ';':
					// rien
				}
			}

			i.rt.ExecPrint(out)
		}
		i.rt.ExecPrint("\n")

	}

	return pc + 1
}

// =======================
// Utils
// =======================

func formatNumber(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%g", f)
}
