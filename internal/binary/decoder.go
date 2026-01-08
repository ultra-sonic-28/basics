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

func DecodeProgram(filename string) (*parser.Program, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// =========================
	// 1️⃣ Lire le HEADER
	// =========================
	var header Header
	if err := binary.Read(f, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	// Magic
	if string(header.Magic[:]) != MagicString {
		return nil, fmt.Errorf("⚠️ invalid magic string")
	}

	// BASIC type
	if _, ok := constants.BasicName[header.BasicType]; !ok {
		return nil, fmt.Errorf("⚠️ unknown BASIC type")
	}

	// Version
	if constants.BasicVersion[header.BasicType] != header.Version {
		return nil, fmt.Errorf("⚠️ BASIC version mismatch")
	}

	// =========================
	// 2️⃣ Lire le payload AST
	// =========================
	astData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// =========================
	// 3️⃣ Vérifier CRC32
	// =========================
	crc := crc32.ChecksumIEEE(astData)
	if crc != header.CRC32 {
		return nil, fmt.Errorf("⚠️ CRC32 mismatch (expected 0x%08X, got 0x%08X)",
			header.CRC32, crc)
	}

	// =========================
	// 4️⃣ Décoder l’AST
	// =========================
	r := bytes.NewReader(astData)
	prog := &parser.Program{}

	for {
		line, err := decodeLine(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		prog.Lines = append(prog.Lines, line)
	}

	// =========================
	// 5️⃣ Infos header
	// =========================
	fmt.Println("📦 BINARY HEADER")
	fmt.Printf("Magic     : %s\n", header.Magic)
	fmt.Printf("Basic type: %s\n", constants.BasicName[header.BasicType])
	fmt.Printf("Version   : %d\n", header.Version)
	fmt.Printf("Nodes     : %d\n", header.NodeCount)
	fmt.Printf("CRC32     : 0x%08X\n\n", header.CRC32)

	return prog, nil
}

func decodeLine(r io.Reader) (*parser.Line, error) {
	lineNum, err := readUint16(r)
	if err != nil {
		return nil, err
	}

	stmtCount, err := readUint16(r)
	if err != nil {
		return nil, err
	}

	line := &parser.Line{
		Number: int(lineNum),
	}

	for i := 0; i < int(stmtCount); i++ {
		stmt, err := decodeStatement(r)
		if err != nil {
			return nil, err
		}
		line.Stmts = append(line.Stmts, stmt)
	}

	return line, nil
}

func decodeStatement(r io.Reader) (parser.Statement, error) {
	op, err := readByte(r)
	if err != nil {
		return nil, err
	}

	switch op {

	case 0x01: // LET
		name, _ := readString(r)
		val, _ := decodeExpression(r)
		return &parser.LetStmt{
			Name:  name,
			Value: val,
		}, nil

	case 0x02: // PRINT
		n, _ := readUint16(r)
		exprs := make([]parser.Expression, 0, n)
		for i := 0; i < int(n); i++ {
			e, _ := decodeExpression(r)
			exprs = append(exprs, e)
		}
		return &parser.PrintStmt{Exprs: exprs}, nil

	case 0x03: // FOR
		name, _ := readString(r)
		start, _ := decodeExpression(r)
		end, _ := decodeExpression(r)
		step, _ := decodeExpression(r)

		return &parser.ForStmt{
			Var:   name,
			Start: start,
			End:   end,
			Step:  step,
		}, nil

	case 0x04: // NEXT
		name, _ := readString(r)
		return &parser.NextStmt{
			Var: name,
		}, nil

	default:
		return nil, fmt.Errorf("decoder: unknown statement opcode 0x%X", op)
	}
}

func decodeExpression(r io.Reader) (parser.Expression, error) {
	op, err := readByte(r)
	if err != nil {
		return nil, err
	}

	switch op {

	case 0x10: // Number
		v, _ := readFloat64(r)
		return &parser.NumberLiteral{Value: v}, nil

	case 0x11: // String
		s, _ := readString(r)
		return &parser.StringLiteral{Value: s}, nil

	case 0x12: // Ident
		name, _ := readString(r)
		return &parser.Identifier{Name: name}, nil

	case 0x13: // Prefix
		opStr, _ := readString(r)
		right, _ := decodeExpression(r)
		return &parser.PrefixExpr{
			Op:    opStr,
			Right: right,
		}, nil

	case 0x14: // Infix
		opStr, _ := readString(r)
		left, _ := decodeExpression(r)
		right, _ := decodeExpression(r)
		return &parser.InfixExpr{
			Op:    opStr,
			Left:  left,
			Right: right,
		}, nil

	default:
		return nil, fmt.Errorf("decoder: unknown expression opcode 0x%X", op)
	}
}

func readByte(r io.Reader) (byte, error) {
	var b byte
	err := binary.Read(r, binary.LittleEndian, &b)
	return b, err
}

func readUint16(r io.Reader) (uint16, error) {
	var v uint16
	err := binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readUint32(r io.Reader) (uint32, error) {
	var v uint32
	err := binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readFloat64(r io.Reader) (float64, error) {
	var v float64
	err := binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readString(r io.Reader) (string, error) {
	l, err := readUint16(r)
	if err != nil {
		return "", err
	}
	buf := make([]byte, l)
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func IsValidBasicsBinary(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var header Header
	if err := binary.Read(f, binary.LittleEndian, &header); err != nil {
		return err
	}

	// Magic
	if string(header.Magic[:]) != MagicString {
		return fmt.Errorf("⚠️ invalid magic string")
	}

	// BASIC type
	if _, ok := constants.BasicName[header.BasicType]; !ok {
		return fmt.Errorf("⚠️ unknown BASIC type")
	}

	// Version
	if constants.BasicVersion[header.BasicType] != header.Version {
		return fmt.Errorf("⚠️ version mismatch")
	}

	return nil
}
