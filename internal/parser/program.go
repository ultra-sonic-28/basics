package parser

func (p *Program) Requires80Columns() bool {
	for _, line := range p.Lines {
		for _, stmt := range line.Stmts {
			if pr, ok := stmt.(*PrStmt); ok {
				if lit, ok := pr.Slot.(*NumberLiteral); ok {
					if int(lit.Value) == 3 {
						return true
					}
				}
			}
		}
	}
	return false
}
