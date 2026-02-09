//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
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
	// Build version
	versionBytes, err := os.ReadFile("VERSION")
	if err != nil {
		return err
	}
	version := strings.TrimSpace(string(versionBytes))

	parts := strings.Split(version, ".")

	if len(parts) != 4 {
		// handle error: unexpected version format
		panic("invalid version string")
	}

	fmt.Println("Building BASICS resource files for version", version)

	tpl := `{
		"RT_GROUP_ICON": {
		  "APP": {
			"0000": [
			  "icon.png",
			  "icon16.png",
			  "icon32.png",
			  "icon48.png",
			  "icon64.png",
			  "icon128.png"
			]
		  }
		},
		"RT_MANIFEST": {
		  "#1": {
			"0409": {
			  "identity": {
				"name": "basics",
				"version": "%[1]s"
			  },
			  "description": "",
			  "minimum-os": "win7",
			  "execution-level": "as invoker",
			  "ui-access": false,
			  "auto-elevate": false,
			  "dpi-awareness": "system",
			  "disable-theming": false,
			  "disable-window-filtering": false,
			  "high-resolution-scrolling-aware": false,
			  "ultra-high-resolution-scrolling-aware": false,
			  "long-path-aware": false,
			  "printer-driver-isolation": false,
			  "gdi-scaling": false,
			  "segment-heap": false,
			  "use-common-controls-v6": false
			}
		  }
		},
		"RT_VERSION": {
		  "#1": {
			"0000": {
			  "fixed": {
				"file_version": "%[1]s",
				"product_version": "%[1]s"
			  },
			  "info": {
				"0409": {
				  "Comments": "",
				  "CompanyName": "ultra-sonic-28",
				  "FileDescription": "BASIC Interpreter for old computers",
				  "FileVersion": "%[1]s",
				  "InternalName": "basics",
				  "LegalCopyright": "© 2025-2026 - ultra-sonic-28 - MIT License",
				  "LegalTrademarks": "",
				  "OriginalFilename": "basics.exe",
				  "PrivateBuild": "",
				  "ProductName": "BASICS",
				  "ProductVersion": "%[1]s",
				  "SpecialBuild": ""
				}
			  }
			}
		  }
		}
	  }
`
	// Générer winres.json
	os.WriteFile(
		"./winres/winres.json",
		[]byte(fmt.Sprintf(tpl, version)),
		0644,
	)

	// Generate windows resource files for embbeding
	cmd := exec.Command("go-winres", "make")
	cmd.Run()

	if err := moveFile("./", "./cmd/basics", "rsrc_windows_386.syso"); err != nil {
		log.Fatal(err)
	}
	if err := moveFile("./", "./cmd/basics", "rsrc_windows_amd64.syso"); err != nil {
		log.Fatal(err)
	}

	// Build binaire
	fmt.Println("Building BASICS binary...")
	builddate := "2025-02-06"
	flags := fmt.Sprintf("-X main.Version=%s -X main.BuildDate=%s", version, builddate)
	cmd = exec.Command(
		"go",
		"build",
		"-ldflags", flags,
		"-o", "./bin/basics.exe",
		"./cmd/basics",
	)
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

// Installing tools : winres
func Tools() error {
	cmd := exec.Command("go", "install", "github.com/tc-hib/go-winres@v0.3.1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Installing tools...")
	return cmd.Run()
}

func moveFile(srcDir, dstDir, name string) error {
	srcPath := filepath.Join(srcDir, name)
	dstPath := filepath.Join(dstDir, name)

	// Ensure destination directory exists.
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	return os.Rename(srcPath, dstPath)
}
