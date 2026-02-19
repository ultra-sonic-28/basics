package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"basics/internal/app"
	"basics/internal/binary"
	"basics/internal/constants"
	"basics/internal/input"
	"basics/internal/interpreter"
	"basics/internal/lexer"
	"basics/internal/logger"
	"basics/internal/machines"
	"basics/internal/parser"
)

var (
	Version   = "dev"
	BuildDate = "unknown"
)

func main() {
	closeLogger, err := logger.InitLogger("basics.log", "basics", logger.LevelInfo)
	if err != nil {
		panic(err)
	}
	defer closeLogger()

	logger.Info("Logging initialized")
	logger.Info("Application starting...")
	logger.Info(fmt.Sprintf("BASICS v%s, built on %s", Version, BuildDate))

	// -------------------------
	// Options CLI
	// -------------------------
	var compileBin bool
	var dumpTokens bool
	var dumpAST bool
	var tty bool
	var basicTypeStr string

	flag.BoolVar(&compileBin, "compile", false, "Generate binary (.bin)")
	flag.BoolVar(&dumpTokens, "dump-tokens", false, "Dump tokens")
	flag.BoolVar(&dumpAST, "dump-ast", false, "Dump AST")
	flag.BoolVar(&tty, "tty", false, "Enable TTY output and ensure that your program does not use any graphical instructions.")
	flag.StringVar(&basicTypeStr, "basic", "APPLE", "BASIC type: APPLE, C64, AMS")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("🆘 Usage: basics [options] <file.bas|file.bin>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	ext := strings.ToLower(filepath.Ext(filename))

	basicType := constants.BASIC_APPLE
	if tty {
		basicType = constants.BASIC_TTY
	}

	// =========================================================
	// Fichier binaire → exécution directe
	// =========================================================
	if ext == ".bin" {

		if compileBin {
			fmt.Println("⚠️ --compile cannot be used with .bin files")
			os.Exit(1)
		}

		// Vérification du header
		if err := binary.IsValidBasicsBinary(filename); err != nil {
			fmt.Println("⚠️ INVALID BINARY PROGRAM")
			os.Exit(1)
		}

		// Décodage binaire → AST
		prog, err := binary.DecodeProgram(filename)
		if err != nil {
			fmt.Printf("⚠️ Error decoding binary: %v\n", err)
			os.Exit(1)
		}

		logger.Info(fmt.Sprintf("Loaded binary file: %s", filename))

		// Exécution
		fmt.Println("\n=== PROGRAM RESULTS ===")
		rt, err := machines.NewRuntime(basicType, false)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		interp := interpreter.New(rt)
		interp.Run(prog)
		return
	}

	// =========================================================
	// Fichier source en BASIC → pipeline classique
	// =========================================================
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("⚠️ Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Loaded source file: %s", filename))

	source := string(data)

	// =========================
	// Lexer
	// =========================
	tokens := lexer.Lex(source)

	if dumpTokens {
		fmt.Println("=== TOKENS ===")
		lexer.DumpTokens(tokens)
	}

	// =========================
	// Parser
	// =========================
	p := parser.New(tokens)
	prog, errs := p.ParseProgram()

	if len(errs) > 0 {
		fmt.Println("\n=== ERRORS ===")
		for _, e := range errs {
			fmt.Println(e.Error())
		}
		os.Exit(1)
	}

	if dumpAST {
		fmt.Println("\n=== AST ===")
		parser.DumpProgram(prog, parser.StdoutEmitter)
	}

	// =========================
	// Compilation binaire
	// =========================
	if compileBin {
		outFile := changeExt(filename, ".bin")

		// déterminer le type BASIC
		var basicType byte
		switch strings.ToUpper(basicTypeStr) {
		case "APPLE":
			basicType = constants.BASIC_APPLE
		case "C64":
			basicType = constants.BASIC_C64
		case "AMS":
			basicType = constants.BASIC_AMS
		default:
			fmt.Printf("Unknown BASIC type '%s', using APPLE\n", basicTypeStr)
			basicType = constants.BASIC_APPLE
		}

		if err := binary.EncodeProgram(prog, outFile, basicType); err != nil {
			fmt.Printf("⚠️ Error during binary compilation: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0) // fin du programme après compilation
	}

	// =========================
	// Interpreter
	// =========================
	if tty {
		fmt.Println("\n=== PROGRAM RESULTS ===")
	}

	mode := prog.Requires80Columns()
	if mode {
		logger.Info("Switch to 80 columns mode")
	} else {
		logger.Info("Switch to 40 columns mode")
	}

	rt, err := machines.NewRuntime(basicType, mode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	interp := interpreter.New(rt)

	// --------------------
	// Mode Terminal (for test and debug purpose)
	// --------------------
	if basicType == constants.BASIC_TTY {
		rt.Input = input.NewTTYInput(os.Stdin, os.Stdout)
		interp.Run(prog)
		return
	}

	// --------------------
	// Mode graphique
	// --------------------
	basicApp := app.NewBasicEbitenApp(rt, interp, prog)
	ebitenApp := app.NewEbitenApp(basicApp)

	if err := ebitenApp.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// changeExt remplace l'extension d'un fichier
func changeExt(path, ext string) string {
	return filepath.Join(filepath.Dir(path),
		filepath.Base(path[:len(path)-len(filepath.Ext(path))])+ext)
}
