package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"

	"basics/internal/constants"
	"basics/internal/parser"
)

// EncodeProgram encode un AST complet dans un fichier binaire avec header + CRC32
func EncodeProgram(prog *parser.Program, filename string, basicType byte) error {
	outFile := filename[:len(filename)-4] + ".bin"

	// =========================
	// Encoder AST en mémoire
	// =========================
	var astBuf bytes.Buffer

	nodeCount := countNodes(prog)

	for _, line := range prog.Lines {
		if err := encodeLine(line, &astBuf); err != nil {
			return err
		}
	}

	// =========================
	// Calcul CRC32 (AST uniquement)
	// =========================
	crc := crc32.ChecksumIEEE(astBuf.Bytes())

	// =========================
	// Construire le header
	// =========================
	var header Header
	copy(header.Magic[:], MagicString)
	header.BasicType = basicType
	header.Version = constants.BasicVersion[basicType]
	header.NodeCount = uint32(nodeCount)
	header.CRC32 = crc

	// =========================
	// Écriture fichier final
	// =========================
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Header
	if err := binary.Write(f, binary.LittleEndian, &header); err != nil {
		return err
	}

	// AST
	if _, err := f.Write(astBuf.Bytes()); err != nil {
		return err
	}

	// =========================
	// Affichage header
	// =========================
	info, err := os.Stat(outFile)
	if err != nil {
		return err
	}

	fmt.Printf("✅ BINARY FILE GENERATED: %s\n", outFile)
	fmt.Printf("Magic      : %s\n", header.Magic)
	fmt.Printf("Basic type : %s\n", constants.BasicName[basicType])
	fmt.Printf("Version    : %d\n", constants.BasicVersion[basicType])
	fmt.Printf("Nodes      : %d\n", nodeCount)
	fmt.Printf("CRC32      : 0x%08X\n", crc)
	fmt.Printf("File size  : %d bytes\n", info.Size())

	return nil
}

// --- les autres fonctions encodeLine / encodeStatement / encodeExpression / writeByte / writeString restent identiques ---

// countNodes parcourt récursivement l'AST pour compter tous les nodes
func countNodes(prog *parser.Program) int {
	count := 0
	for _, line := range prog.Lines {
		count++ // line
		for _, stmt := range line.Stmts {
			if stmt != nil {
				count += countStatementNodes(stmt)
			}
		}
	}
	return count
}

func countStatementNodes(stmt parser.Statement) int {
	if stmt == nil {
		return 0 // REM
	}

	count := 1

	switch s := stmt.(type) {
	case *parser.LetStmt:
		count += countExprNodes(s.Value)

	case *parser.PrintStmt:
		for _, e := range s.Exprs {
			count += countExprNodes(e)
		}

	case *parser.ForStmt:
		count += countExprNodes(s.Start)
		count += countExprNodes(s.End)
		if s.Step != nil {
			count += countExprNodes(s.Step)
		}

	case *parser.NextStmt:
		// rien
	}

	return count
}

func countExprNodes(expr parser.Expression) int {
	switch e := expr.(type) {
	case *parser.NumberLiteral, *parser.StringLiteral, *parser.Identifier:
		return 1
	case *parser.PrefixExpr:
		return 1 + countExprNodes(e.Right)
	case *parser.InfixExpr:
		return 1 + countExprNodes(e.Left) + countExprNodes(e.Right)
	default:
		return 0
	}
}

func encodeLine(line *parser.Line, w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint16(line.Number)); err != nil {
		return err
	}

	// Compter uniquement les statements non-nil
	count := uint16(0)
	for _, stmt := range line.Stmts {
		if stmt != nil {
			count++
		}
	}

	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}

	for _, stmt := range line.Stmts {
		if stmt == nil {
			// REM → ignoré
			continue
		}
		if err := encodeStatement(stmt, w); err != nil {
			return err
		}
	}

	return nil
}

func encodeStatement(stmt parser.Statement, w io.Writer) error {
	if stmt == nil {
		// REM → ignoré
		return nil
	}

	switch s := stmt.(type) {
	case *parser.LetStmt:
		if err := writeByte(w, 0x01); err != nil {
			return err
		}
		if err := writeString(w, s.Name); err != nil {
			return err
		}
		if err := encodeExpression(s.Value, w); err != nil {
			return err
		}

	case *parser.PrintStmt:
		if err := writeByte(w, 0x02); err != nil {
			return err
		}
		exprCount := uint16(len(s.Exprs))
		if err := binary.Write(w, binary.LittleEndian, exprCount); err != nil {
			return err
		}
		for _, e := range s.Exprs {
			if err := encodeExpression(e, w); err != nil {
				return err
			}
		}

	case *parser.ForStmt:
		if err := writeByte(w, 0x03); err != nil {
			return err
		}
		if err := writeString(w, s.Var); err != nil {
			return err
		}
		if err := encodeExpression(s.Start, w); err != nil {
			return err
		}
		if err := encodeExpression(s.End, w); err != nil {
			return err
		}
		step := s.Step
		if step == nil {
			step = &parser.NumberLiteral{Value: 1}
		}
		if err := encodeExpression(step, w); err != nil {
			return err
		}

	case *parser.NextStmt:
		if err := writeByte(w, 0x04); err != nil {
			return err
		}
		if err := writeString(w, s.Var); err != nil {
			return err
		}

	default:
		return fmt.Errorf("encoder: statement not supported %T", stmt)
	}
	return nil
}

func encodeExpression(expr parser.Expression, w io.Writer) error {
	switch e := expr.(type) {
	case *parser.NumberLiteral:
		if err := writeByte(w, 0x10); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, e.Value); err != nil {
			return err
		}
	case *parser.StringLiteral:
		if err := writeByte(w, 0x11); err != nil {
			return err
		}
		if err := writeString(w, e.Value); err != nil {
			return err
		}
	case *parser.Identifier:
		if err := writeByte(w, 0x12); err != nil {
			return err
		}
		if err := writeString(w, e.Name); err != nil {
			return err
		}
	case *parser.PrefixExpr:
		if err := writeByte(w, 0x13); err != nil {
			return err
		}
		if err := writeString(w, e.Op); err != nil {
			return err
		}
		if err := encodeExpression(e.Right, w); err != nil {
			return err
		}
	case *parser.InfixExpr:
		if err := writeByte(w, 0x14); err != nil {
			return err
		}
		if err := writeString(w, e.Op); err != nil {
			return err
		}
		if err := encodeExpression(e.Left, w); err != nil {
			return err
		}
		if err := encodeExpression(e.Right, w); err != nil {
			return err
		}
	default:
		return fmt.Errorf("encoder: expression not supported %T", expr)
	}
	return nil
}

func writeByte(w io.Writer, b byte) error {
	return binary.Write(w, binary.LittleEndian, b)
}

func writeString(w io.Writer, s string) error {
	l := uint16(len(s))
	if err := binary.Write(w, binary.LittleEndian, l); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}
