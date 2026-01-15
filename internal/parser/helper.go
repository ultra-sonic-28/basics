package parser

func StmtName(s Statement) string {
	switch s.(type) {
	case *HomeStmt:
		return "HOME"
	case *PrintStmt:
		return "PRINT"
	case *LetStmt:
		return "LET"
	case *IfStmt:
		return "IF"
	case *IfJumpStmt:
		return "IFMULTI"
	case *GotoStmt:
		return "GOTO"
	case *GosubStmt:
		return "GOSUB"
	case *ReturnStmt:
		return "RETURN"
	case *ForStmt:
		return "FOR"
	case *NextStmt:
		return "NEXT"
	case *EndStmt:
		return "END"
	case *HTabStmt:
		return "HTAB"
	case *VTabStmt:
		return "VTAB"
	default:
		return "UNKNOWN"
	}
}
