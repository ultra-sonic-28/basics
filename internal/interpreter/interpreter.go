package interpreter

import (
	"fmt"
	"strings"

	"basics/internal/errors"
	"basics/internal/parser"
	"basics/internal/runtime"
)

// ForFrame garde les infos d'une boucle FOR active
type ForFrame struct {
	Var       string
	End       float64
	Step      float64
	LineIndex int // index dans prog.Lines
	StmtIndex int // index du statement FOR dans la ligne
}

// ForStack gère une pile de boucles FOR
type ForStack struct {
	stack []ForFrame
}

func NewForStack() *ForStack {
	return &ForStack{}
}

func (fs *ForStack) Push(f ForFrame) {
	fs.stack = append(fs.stack, f)
}

func (fs *ForStack) Pop() ForFrame {
	n := len(fs.stack)
	f := fs.stack[n-1]
	fs.stack = fs.stack[:n-1]
	return f
}

func (fs *ForStack) Top() *ForFrame {
	if len(fs.stack) == 0 {
		return nil
	}
	return &fs.stack[len(fs.stack)-1]
}

// Interpreter exécute un programme BASIC
type Interpreter struct {
	rt    *runtime.Runtime
	forSt *ForStack
}

func New(rt *runtime.Runtime) *Interpreter {
	return &Interpreter{
		rt:    rt,
		forSt: NewForStack(),
	}
}

// Run exécute le programme
func (i *Interpreter) Run(prog *parser.Program) {
	lines := prog.Lines

	for pc := 0; pc < len(lines); pc++ {
		line := lines[pc]

		// stmtIdx = index du statement courant dans la ligne
		for stmtIdx := 0; stmtIdx < len(line.Stmts); stmtIdx++ {
			stmt := line.Stmts[stmtIdx]

			switch s := stmt.(type) {

			case *parser.LetStmt:
				val, err := EvalExpr(s.Value, i.rt)
				if err != nil {
					fmt.Println(err)
					return
				}
				i.rt.Env.Set(s.Name, val)

			case *parser.PrintStmt:
				cursor := 0 // position "colonne" simulée pour la virgule
				for iExpr, expr := range s.Exprs {
					val, err := EvalExpr(expr, i.rt)
					if err != nil {
						fmt.Println(err)
						return
					}

					// afficher la valeur
					str := ""
					switch val.Type {
					case runtime.NUMBER:
						str = formatNumber(val.Num)
					case runtime.STRING:
						str = val.Str
					}

					// si ce n'est pas la première expression, gérer le séparateur
					if iExpr > 0 {
						// séparateur entre deux expressions dans PRINT
						sep := s.Separators[iExpr-1] // tableau de séparateurs ';' ou ',' défini lors du parsing
						switch sep {
						case ';':
							// rien, juste continuer après la précédente valeur
						case ',':
							// avancé jusqu'à la prochaine tabulation de 14 colonnes
							spaces := 14 - (cursor % 14)
							str = strings.Repeat(" ", spaces) + str
						}
					}

					i.rt.ExecPrint(str)
					cursor += len(str)
				}
				i.rt.ExecPrint("\n")

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
							s.LineNum,
							"STEP CANNOT BE ZERO",
						))
						return
					}
				}

				i.rt.Env.Set(s.Var, runtime.Value{Type: runtime.NUMBER, Num: startVal.Num})
				i.forSt.Push(ForFrame{
					Var:       s.Var,
					End:       endVal.Num,
					Step:      step,
					LineIndex: pc,
					StmtIndex: stmtIdx,
				})

			case *parser.NextStmt:
				frame := i.forSt.Top()
				if frame == nil {
					fmt.Printf("?NEXT WITHOUT FOR\n")
					return
				}

				v, _ := i.rt.Env.Get(frame.Var)
				v.Num += frame.Step

				// Condition selon le signe du step
				done := (frame.Step > 0 && v.Num > frame.End) || (frame.Step < 0 && v.Num < frame.End)
				if !done {
					i.rt.Env.Set(frame.Var, v)
					pc = frame.LineIndex
				} else {
					i.forSt.Pop()
				}

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
			}
		}
	}
}

func formatNumber(f float64) string {
	if f == float64(int64(f)) {
		// entier → pas de décimales
		return fmt.Sprintf("%d", int64(f))
	}
	// réel → afficher avec décimales
	return fmt.Sprintf("%g", f)
}
