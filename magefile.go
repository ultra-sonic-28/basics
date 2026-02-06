//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	mainCmd   = ".\\cmd\\basic\\main.go"
	basicsDir = "basics"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Run

// Build basics binary to /bin directory
func Build() error {
	fmt.Println("Building basics binary...")
	cmd := exec.Command("go", "build", "-o", "./bin/basics.exe", "./cmd/basics")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run unit tests with coverage support
func Test() error {
	cmd := exec.Command("go", "run", ".\\test_summary.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running tests...")
	return cmd.Run()
}

// Delete all .exe files under basics/ except in folders starting with "."
func Clean() error {
	fmt.Println("Cleaning executables in", basicsDir)

	return filepath.Walk(basicsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore dossiers cachés
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Supprime les .exe
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".exe") {
			fmt.Println("Removing:", path)
			return os.Remove(path)
		}

		return nil
	})
}
