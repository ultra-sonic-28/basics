package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"basics/internal/common"
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

		// index de la première instruction de la ligne
		if _, exists := i.lineIndex[line.Number]; !exists {
			i.lineIndex[line.Number] = len(i.insts)
		}

		for _, stmt := range line.Stmts {

			switch s := stmt.(type) {

			// =====================================================
			// IF : aplatissement en flot linéaire (style Applesoft)
			// =====================================================
			case *parser.IfStmt:

				// 1️⃣ réserver une instruction IF (patchée ensuite)
				ifPC := len(i.insts)

				i.insts = append(i.insts, Instruction{
					LineNum: line.Number,
					Stmt:    nil, // sera remplacé
				})

				// 2️⃣ THEN block (instructions normales)
				for _, thenStmt := range s.Then {
					i.insts = append(i.insts, Instruction{
						LineNum: line.Number,
						Stmt:    thenStmt,
					})
				}

				// 3️⃣ ELSE block (optionnel)
				var elseTarget int
				if len(s.Else) > 0 {

					// saut après THEN
					gotoAfterThenPC := len(i.insts)

					i.insts = append(i.insts, Instruction{
						LineNum: line.Number,
						Stmt: &parser.GotoStmt{
							Expr: &parser.NumberLiteral{
								Value: float64(-1), // patch plus tard
							},
						},
					})

					elseTarget = len(i.insts)

					for _, elseStmt := range s.Else {
						i.insts = append(i.insts, Instruction{
							LineNum: line.Number,
							Stmt:    elseStmt,
						})
					}

					// patch du GOTO de fin de THEN
					i.insts[gotoAfterThenPC].Stmt.(*parser.GotoStmt).Expr =
						&parser.NumberLiteral{Value: float64(len(i.insts))}

				} else {
					elseTarget = len(i.insts)
				}

				// 4️⃣ patch de l’instruction IF
				i.insts[ifPC].Stmt = &parser.IfJumpStmt{
					Cond:   s.Cond,
					Target: elseTarget,
				}

			// ==========================
			// Autres instructions
			// ==========================
			default:
				i.insts = append(i.insts, Instruction{
					LineNum: line.Number,
					Stmt:    stmt,
				})
			}
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
	logger.Debug("Program execution trace")
	logger.Debug(fmt.Sprintf("Program contains %d lines and %d instructions", len(prog.Lines), len(i.insts)))

	pc := 0
	for pc < len(i.insts) {
		inst := i.insts[pc]
		nextPC := pc + 1
		sExpr := ""

		switch s := inst.Stmt.(type) {

		// -----------------------
		// PR
		// -----------------------
		case *parser.PrStmt:
			val, _, err := EvalExpr(s.Slot, i.rt)
			if err != nil {
				fmt.Println(err)
				return
			}
			sExpr = fmt.Sprintf("-> %g", val.Num)
			i.rt.ExecPr(int(val.Num))

		// -----------------------
		// REM
		// -----------------------
		case *parser.RemStmt:
			sExpr = s.Text

		// -----------------------
		// HOME
		// -----------------------
		case *parser.HomeStmt:
			i.rt.ExecHome()

		// -----------------------
		// END
		// -----------------------
		case *parser.EndStmt:
			logger.Debug(LogTrace(inst, pc, nextPC, sExpr))
			i.rt.Halt()
			return

		// -----------------------
		// LET
		// -----------------------
		case *parser.LetStmt:
			sExpr = i.execLet(s, &inst)

		// -----------------------
		// INPUT
		// -----------------------
		case *parser.InputStmt:
			i.execInput(s)

		// -----------------------
		// GET
		// -----------------------
		case *parser.GetStmt:
			i.execGet(s)
			val, _ := i.rt.Env.Get(s.Var.Name)
			sExpr = val.Str

		// -----------------------
		// PRINT
		// -----------------------
		case *parser.PrintStmt:
			// PRINT sans arguments
			if len(s.Exprs) == 0 {
				i.rt.ExecPrint("\n")
				break
			}

			cursor := i.rt.Video.CursorX()

			for iExpr, expr := range s.Exprs {

				// 1. Gérer le séparateur venant de l'expression PRÉCÉDENTE
				if iExpr > 0 {
					sep := s.Separators[iExpr-1]
					if sep == ',' {
						spaces := 16 - (cursor % 16)
						if spaces == 0 {
							spaces = 16
						}
						padding := strings.Repeat(" ", spaces)
						i.rt.ExecPrint(padding)
						cursor = i.rt.Video.CursorX()
					}
				}

				// ==========================
				// GESTION DE TAB
				// ==========================
				if tabExpr, ok := expr.(*parser.TabExpr); ok {

					val, _, err := EvalExpr(tabExpr.Expr, i.rt)
					if err != nil {
						fmt.Println(err)
						return
					}

					if val.Type == runtime.STRING {
						i.rt.ExecError(errors.NewSemantic(
							tabExpr.Line,
							"EXPECTED NUMBER",
						))
						return
					}

					var target int

					switch val.Type {
					case runtime.INTEGER:
						target = val.Int
					case runtime.NUMBER:
						target = int(val.Num)
					}

					// BASIC TAB is 1-based: TAB(1) is the first column (index 0)
					targetColumn := target - 1
					if targetColumn < 0 {
						targetColumn = 0
					}

					currentX := cursor

					if targetColumn > currentX {
						spaces := targetColumn - currentX
						padding := strings.Repeat(" ", spaces)
						i.rt.ExecPrint(padding)
						cursor = i.rt.Video.CursorX()
					}

					continue
				}

				// ==========================
				// GESTION DE SPC
				// ==========================
				if spcExpr, ok := expr.(*parser.SpcExpr); ok {

					val, _, err := EvalExpr(spcExpr.Expr, i.rt)
					if err != nil {
						fmt.Println(err)
						return
					}

					if val.Type == runtime.STRING {
						i.rt.ExecError(errors.NewSemantic(
							spcExpr.Line,
							"EXPECTED NUMBER",
						))
						return
					}

					var count int

					switch val.Type {
					case runtime.INTEGER:
						count = val.Int
					case runtime.NUMBER:
						count = int(val.Num)
					}

					if count > 0 {
						spaces := strings.Repeat(" ", count)
						i.rt.ExecPrint(spaces)
						cursor = i.rt.Video.CursorX()
					}

					continue
				}

				val, _, err := EvalExpr(expr, i.rt)
				if err != nil {
					i.rt.ExecError(err)
					return
				}

				str := ""
				switch val.Type {
				case runtime.INTEGER:
					str = common.FormatNumber(float64(val.Int))
				case runtime.NUMBER:
					str = common.FormatNumber(val.Num)
				case runtime.STRING:
					str = val.Str
				}

				/* if indices != nil {
					sExpr += fmt.Sprint(common.Flatten(*indices))
				} */
				sExpr += str

				i.rt.ExecPrint(str)
				cursor = i.rt.Video.CursorX()
			}

			if len(s.Separators) < len(s.Exprs) {
				i.rt.ExecPrint("\n")
			}

		// -----------------------
		// NORMAL / INVERSE
		// -----------------------
		case *parser.NormalStmt:
			i.rt.SetInverse(false)
			i.rt.SetFlash(false)

		case *parser.InverseStmt:
			i.rt.SetInverse(true)
			i.rt.SetFlash(false)

		// -----------------------
		// FLASH
		// -----------------------
		case *parser.FlashStmt:
			i.rt.SetInverse(false)
			i.rt.SetFlash(true)

		// -----------------------
		// CLEAR
		// -----------------------
		case *parser.ClearStmt:
			i.rt.Clear()

		// -----------------------
		// DIM
		// -----------------------
		case *parser.DimStmt:
			sExpr = i.execDim(s)

		// -----------------------
		// HTAB / VTAB
		// -----------------------
		case *parser.HTabStmt:
			val, _, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return
			}

			sExpr = fmt.Sprintf("%d", int(val.Num))
			i.rt.ExecHTab(int(val.Num))

		case *parser.VTabStmt:
			val, _, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return
			}

			sExpr = fmt.Sprintf("%d", int(val.Num))
			i.rt.ExecVTab(int(val.Num))

		// -----------------------
		// FOR (Applesoft semantics)
		// -----------------------
		case *parser.ForStmt:
			startVal, _, err := EvalExpr(s.Start, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return
			}

			endVal, _, err := EvalExpr(s.End, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return
			}

			step := 1.0
			if s.Step != nil {
				stepVal, _, err := EvalExpr(s.Step, i.rt)
				if err != nil {
					fmt.Println(err)
					return
				}
				step = stepVal.Num
				if step == 0 {
					err = errors.NewSemantic(
						inst.LineNum,
						"STEP CANNOT BE ZERO",
					)
					i.rt.ExecError(err)
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

			sExpr = fmt.Sprintf("-> %g TO %g STEP %g", startVal.Num, endVal.Num, step)

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
				sExpr = fmt.Sprintf("-> %g", v.Num)
			} else {
				i.forStack.Pop()
			}

		// -----------------------
		// GOTO
		// -----------------------
		case *parser.GotoStmt:
			val, _, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return
			}

			if val.Type != runtime.NUMBER {
				fmt.Println("?GOTO TYPE MISMATCH")
				return
			}

			line := int(val.Num)
			sExpr = fmt.Sprintf("%d", line)
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
			val, _, err := EvalExpr(s.Expr, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return
			}

			if val.Type != runtime.NUMBER {
				fmt.Println("?GOSUB TYPE MISMATCH")
				return
			}

			line := int(val.Num)
			sExpr = fmt.Sprintf("%d", line)
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
			cond, _, err := EvalExpr(s.Cond, i.rt)
			if err != nil {
				i.rt.ExecError(err)
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
				// exécution inline TERMINALE
				pc2 := pc + 1 // PC logique après le IF

				for _, stmt := range s.Then {
					pc2 = i.execInline(inst.LineNum, stmt, pc2-1)
				}

				nextPC = pc2
				sExpr = "THEN"
			} else if s.Else != nil {
				pc2 := pc + 1
				for _, stmt := range s.Else {
					pc2 = i.execInline(inst.LineNum, stmt, pc2-1)
				}
				nextPC = pc2
				sExpr = "ELSE"
			} else {
				// condition fausse → instruction suivante
				nextPC = pc + 1
				sExpr = "ELSE"
			}

		// -----------------------
		// IF (compiled jump)
		// -----------------------
		case *parser.IfJumpStmt:
			cond, _, err := EvalExpr(s.Cond, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return
			}

			exec := false
			switch cond.Type {
			case runtime.BOOLEAN:
				exec = cond.Flag
			case runtime.NUMBER:
				exec = cond.Num != 0
			}

			sExpr = "THEN"
			if !exec {
				nextPC = s.Target
				sExpr = "ELSE"
			}

		}

		logger.Debug(LogTrace(inst, pc, nextPC, sExpr))
		pc = nextPC
	}

	i.rt.Video.Render()
}

// =======================
// Inline execution helper
// =======================

func (i *Interpreter) execInline(line int, stmt parser.Statement, pc int) int {
	_ = line

	switch s := stmt.(type) {

	case *parser.HomeStmt:
		i.rt.ExecHome()
		return pc + 1

	case *parser.GotoStmt:
		val, _, err := EvalExpr(s.Expr, i.rt)
		if err != nil {
			i.rt.ExecError(err)
			return pc + 1
		}
		target, ok := i.lineIndex[int(val.Num)]
		if !ok {
			fmt.Printf("?UNDEFINED LINE %d\n", int(val.Num))
			return pc + 1
		}
		return target

	case *parser.GosubStmt:
		val, _, err := EvalExpr(s.Expr, i.rt)
		if err != nil {
			i.rt.ExecError(err)
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
		val, _, err := EvalExpr(s.Value, i.rt)
		if err != nil {
			i.rt.ExecError(err)
			return pc + 1
		}
		i.rt.Env.Set(s.Name, val)
		return pc + 1

	case *parser.GetStmt:
		i.execGet(s)
		return pc + 1

	case *parser.DimStmt:
		_ = i.execDim(s)
		return pc + 1

	case *parser.PrintStmt:
		//i.rt.ExecPrint(s.Exprs[0].(*parser.StringLiteral).Value)
		//i.rt.ExecPrint("\n")
		cursor := i.rt.Video.CursorX()
		for iExpr, expr := range s.Exprs {

			if iExpr > 0 {
				sep := s.Separators[iExpr-1]
				if sep == ',' {
					spaces := 16 - (cursor % 16)
					if spaces == 0 {
						spaces = 16
					}
					padding := strings.Repeat(" ", spaces)
					i.rt.ExecPrint(padding)
					cursor = i.rt.Video.CursorX()
				}
			}

			// GESTION DE TAB/SPC DANS INLINE
			if tabExpr, ok := expr.(*parser.TabExpr); ok {
				val, _, _ := EvalExpr(tabExpr.Expr, i.rt)

				var target int
				switch val.Type {
				case runtime.INTEGER:
					target = val.Int
				case runtime.NUMBER:
					target = int(val.Num)
				}

				// BASIC TAB is 1-based: TAB(1) is the first column (index 0)
				targetColumn := target - 1
				if targetColumn < 0 {
					targetColumn = 0
				}

				currentX := cursor
				if targetColumn > currentX {
					spaces := targetColumn - currentX
					i.rt.ExecPrint(strings.Repeat(" ", spaces))
					cursor = i.rt.Video.CursorX()
				}
				continue
			}

			if spcExpr, ok := expr.(*parser.SpcExpr); ok {
				val, _, _ := EvalExpr(spcExpr.Expr, i.rt)

				var count int
				switch val.Type {
				case runtime.INTEGER:
					count = val.Int
				case runtime.NUMBER:
					count = int(val.Num)
				}

				if count > 0 {
					i.rt.ExecPrint(strings.Repeat(" ", count))
					cursor = i.rt.Video.CursorX()
				}
				continue
			}

			val, _, err := EvalExpr(expr, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return pc
			}

			var out string
			switch val.Type {
			case runtime.STRING:
				out = val.Str
			case runtime.NUMBER:
				out = common.FormatNumber(val.Num)
			case runtime.BOOLEAN:
				if val.Flag {
					out = "1"
				} else {
					out = "0"
				}
			case runtime.INTEGER:
				out = strconv.Itoa(val.Int)
			default:
				out = ""
			}

			i.rt.ExecPrint(out)
			cursor = i.rt.Video.CursorX()
		}

		if len(s.Separators) < len(s.Exprs) {
			i.rt.ExecPrint("\n")
		}

	}

	return pc + 1
}

func (i *Interpreter) execLet(s *parser.LetStmt, inst *Instruction) string {
	sExpr := ""

	val, _, err := EvalExpr(s.Value, i.rt)
	if err != nil {
		i.rt.ExecError(err)
		return ""
	}

	// ======================================================
	// CAS 1 : VARIABLE SIMPLE (comportement existant)
	// ======================================================
	if len(s.Indices) == 0 {

		vType := common.VarType(s.Name)
		switch vType {

		case "int":
			if val.Type == runtime.STRING {
				err := errors.NewSemantic(inst.LineNum, "TYPE MISMATCH: INTEGER EXPECTED")
				i.rt.ExecError(err)
				return ""
			}

			// Coercition : troncature vers entier
			intValue := int(val.Num)
			if val.Type == runtime.INTEGER {
				intValue = val.Int
			}

			i.rt.Env.Set(s.Name, runtime.Value{
				Type: runtime.INTEGER,
				Int:  intValue,
			})
			sExpr = fmt.Sprintf("%d", intValue)

		case "string":
			if val.Type != runtime.STRING {
				err := errors.NewSemantic(inst.LineNum, "TYPE MISMATCH: STRING EXPECTED")
				i.rt.ExecError(err)
				return ""
			}
			i.rt.Env.Set(s.Name, val)
			sExpr = val.Str

		case "float":
			var v float64

			if val.Type == runtime.STRING {
				err := errors.NewSemantic(inst.LineNum, "TYPE MISMATCH: FLOAT EXPECTED")
				i.rt.ExecError(err)
				return ""
			}

			if val.Type == runtime.INTEGER {
				v = float64(val.Int)
			} else if val.Type == runtime.NUMBER {
				v = val.Num
			}

			i.rt.Env.Set(s.Name, runtime.Value{
				Type: runtime.NUMBER,
				Num:  v,
			})
			sExpr = fmt.Sprintf("%g", v)

		}

		return sExpr
	}

	// ======================================================
	// CAS 2 : AFFECTATION DANS UN TABLEAU
	// ======================================================
	arrVal, ok := i.rt.Env.Get(s.Name)
	if !ok || arrVal.Type != runtime.ARRAY {
		err := errors.NewSemantic(inst.LineNum, "BAD SUBSCRIPT")
		i.rt.ExecError(err)
		return ""
	}

	arr := arrVal.Array

	// --- évaluation des indices ---
	var indices []int
	for _, idxExpr := range s.Indices {
		v, _, err := EvalExpr(idxExpr, i.rt)
		if err != nil {
			i.rt.ExecError(err)
			return ""
		}

		if v.Type != runtime.NUMBER && v.Type != runtime.INTEGER {
			err := errors.NewSemantic(inst.LineNum, "BAD SUBSCRIPT")
			i.rt.ExecError(err)
			return ""
		}

		indices = append(indices, int(v.Num))
	}

	// --- vérification du type stocké ---
	switch arr.BaseType {

	case runtime.INTEGER:
		if val.Type == runtime.STRING {
			err := errors.NewSemantic(inst.LineNum, "TYPE MISMATCH: INTEGER EXPECTED")
			i.rt.ExecError(err)
			return ""
		}

		intValue := int(val.Num)
		if val.Type == runtime.INTEGER {
			intValue = val.Int
		}
		val = runtime.Value{Type: runtime.INTEGER, Int: intValue}

	case runtime.STRING:
		if val.Type != runtime.STRING {
			err := errors.NewSemantic(inst.LineNum, "TYPE MISMATCH: STRING EXPECTED")
			i.rt.ExecError(err)
			return ""
		}

	case runtime.NUMBER:
		if val.Type == runtime.STRING {
			err := errors.NewSemantic(inst.LineNum, "TYPE MISMATCH: FLOAT EXPECTED")
			i.rt.ExecError(err)
			return ""
		}
		val = runtime.Value{Type: runtime.NUMBER, Num: val.Num}
	}

	// --- affectation ---
	if err := arr.Set(indices, val); err != nil {
		i.rt.ExecError(errors.NewSemantic(inst.LineNum, err.Error()))
		return ""
	}

	sExpr = val.String()
	return sExpr
}

func (i *Interpreter) execDim(s *parser.DimStmt) string {
	var allVars strings.Builder
	allVars.WriteString("-> (")
	for _, decl := range s.Arrays {

		var dims []int
		for j, expr := range decl.Dimensions {
			v, _, err := EvalExpr(expr, i.rt)
			if err != nil {
				i.rt.ExecError(err)
				return ""
			}

			if v.Type != runtime.NUMBER && v.Type != runtime.INTEGER {
				err = errors.NewSemantic(
					s.Line,
					"BAD SUBSCRIPT: DIMENSION MUST BE A NUMBER OR AN INTEGER",
				)
				i.rt.ExecError(err)
				return ""
			}

			size := int(v.Num) + 1
			if size <= 0 {
				err = errors.NewSemantic(
					s.Line,
					"BAD SUBSCRIPT: DIMENSION MUST BE POSITIVE",
				)
				i.rt.ExecError(err)
				return ""
			}

			if j > 0 {
				allVars.WriteString(", ")
			}
			allVars.WriteString(fmt.Sprintf("%d", int(v.Num)))
			dims = append(dims, size)
		}

		arr := runtime.NewArray(decl.BaseType, dims)

		i.rt.Env.Set(decl.Name, runtime.Value{
			Type:  runtime.ARRAY,
			Array: arr,
		})
	}

	allVars.WriteString(")")
	return allVars.String()
}

func (i *Interpreter) execInput(s *parser.InputStmt) {
	for {
		// afficher le prompt
		if s.Prompt != nil {
			i.rt.ExecPrint(s.Prompt.Value)
		} else {
			i.rt.ExecPrint("? ")
		}

		line, _ := i.rt.ExecInput()
		line = strings.TrimRight(line, "\r\n")

		values := strings.Split(line, ",")

		if len(values) != len(s.Vars) {
			i.rt.ExecPrint("\n?NOT ENOUGH VALUES, REENTER\n")
			continue
		}

		ok := true

		for idx, v := range s.Vars {
			val := strings.TrimSpace(values[idx])

			// STRING variable
			if strings.HasSuffix(v.Name, "$") {
				i.rt.Env.Set(v.Name, runtime.Value{
					Type: runtime.STRING,
					Str:  strings.TrimSpace(val),
				})
				continue
			}

			// NUMERIC variable
			numStr := strings.ReplaceAll(val, " ", "")
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				ok = false
				break
			}

			// INTEGER or NUMBER
			if strings.HasSuffix(v.Name, "%") {
				i.rt.Env.Set(v.Name, runtime.Value{
					Type: runtime.INTEGER,
					Int:  int(num),
				})
			} else {
				i.rt.Env.Set(v.Name, runtime.Value{
					Type: runtime.NUMBER,
					Num:  num,
				})
			}
		}

		if !ok {
			i.rt.ExecPrint("\n?TYPE MISMATCH, REENTER\n")
			continue
		}

		//i.rt.ExecPrint("\n")
		break
	}

	i.rt.DisableKeyboard()
}

func (i *Interpreter) execGet(s *parser.GetStmt) {
	// lecture bloquante d'un caractère
	for {
		ch, err := i.rt.ExecGet()
		if err != nil {
			i.rt.ExecError(err)
			return
		}

		if strings.HasSuffix(s.Var.Name, "$") {
			i.rt.Env.Set(s.Var.Name, runtime.Value{
				Type: runtime.STRING,
				Str:  string(ch),
			})

			break
		}

		val := strings.TrimSpace(string(ch))

		// NUMERIC variable
		numStr := strings.ReplaceAll(val, " ", "")
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			i.rt.ExecPrint("\n?TYPE MISMATCH, REENTER\n")
			continue
		}

		// INTEGER or NUMBER
		if strings.HasSuffix(s.Var.Name, "%") {
			i.rt.Env.Set(s.Var.Name, runtime.Value{
				Type: runtime.INTEGER,
				Int:  int(num),
			})
		} else {
			i.rt.Env.Set(s.Var.Name, runtime.Value{
				Type: runtime.NUMBER,
				Num:  num,
			})
		}

		break
	}
}

// =======================
// Utils
// =======================

func LogTrace(inst Instruction, pc int, nextPC int, sExpr string) string {
	return fmt.Sprintf(
		"Executing line: %d, pc: %d, nextPC: %d - [%s]%s %s",
		inst.LineNum,
		pc,
		nextPC,
		parser.StmtName(inst.Stmt),
		parser.StmtArgs(inst.Stmt),
		sExpr,
	)
}
