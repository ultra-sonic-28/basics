//go:build mage
// +build mage

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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
	// Incrémenter le numéro de build AVANT toute lecture
	fmt.Println("Updating build number")
	versionBytes, err := incrementBuildNumber("VERSION")
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

// Running sources analysis and statistics generation
func Stats() error {
	cloc_executable := "./.bintools/cloc-2.06.exe"
	report_file := "./.tmp/cloc_output_raw.md"

	cmd := exec.Command(
		cloc_executable,
		"--skip-uniqueness",
		"--quiet",
		"--skip-archive=(zip|tar(.(gz|Z|bz2|xz|7z))?)",
		"--skip-win-hidden",
		//"--thousands-delimiter=_",
		//"--fmt=2",
		"--md",
		"--report-file="+report_file,
		"--found=./.tmp/found.txt",
		"--ignored=./.tmp/ignored.txt",
		"./cmd/",
		"./examples/",
		"./internal/",
		"./testutils/",
		"./winres/",
		"./architecture.md",
		"./CHANGELOG.md",
		"./README.md",
		"./*.go",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running stats...")
	err := cmd.Run()

	fmt.Println("Formating stats...")
	formatStatsFile()

	return err
}

// Utility functions

func moveFile(srcDir, dstDir, name string) error {
	srcPath := filepath.Join(srcDir, name)
	dstPath := filepath.Join(dstDir, name)

	// Ensure destination directory exists.
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	return os.Rename(srcPath, dstPath)
}

// formatThousands insère _ tous les trois chiffres à partir de la droite.
func formatThousands(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	out := []byte{}
	count := 0
	for i := n - 1; i >= 0; i-- {
		out = append(out, s[i])
		count++
		if count%3 == 0 && i != 0 {
			out = append(out, '_')
		}
	}
	// reverse
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

// Capitalise le premier caractère de chaque "mot" entre les pipes.
func capitalizeHeader(line string) string {
	parts := strings.Split(line, "|")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			runes := []rune(p)
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "|")
}

func formatStatsFile() {
	input := "./.tmp/cloc_output_raw.md"
	output := "./.tmp/cloc_output_clean.md"

	in, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	lineNum := 0
	reNumber := regexp.MustCompile(`\b\d{1,3}(\d{3})*\b`)

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Ignore les trois premières lignes du fichier
		if lineNum <= 3 {
			continue
		}

		// Supprime la ligne de tirets "--------|--------|..."
		if strings.HasPrefix(line, "--------|") {
			continue
		}

		// Capitaliser l'en-tête Language|files|blank|comment|code
		if strings.HasPrefix(line, "Language|") {
			line = capitalizeHeader(line)
		}

		// Substitutions de texte
		line = strings.ReplaceAll(line, "Visual Basic", "Basic")
		line = strings.ReplaceAll(line, "SUM:", "TOTAL:")

		// Formattage des nombres avec _
		line = reNumber.ReplaceAllStringFunc(line, func(num string) string {
			return formatThousands(num)
		})

		// Si la ligne contient "TOTAL:", mettre chaque cellule de TOTAL en gras
		if strings.HasPrefix(line, "TOTAL:|") || strings.HasPrefix(line, "**TOTAL:**|") {
			parts := strings.Split(line, "|")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
				if parts[i] != "" {
					parts[i] = fmt.Sprintf("**%s**", parts[i])
				}
			}
			line = strings.Join(parts, "|")
		}

		fmt.Fprintln(writer, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	writer.Flush()
}

func incrementBuildNumber(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(data))
	parts := strings.Split(version, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid version string: %q", version)
	}

	build, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", fmt.Errorf("invalid build number %q: %w", parts[3], err)
	}
	build++

	parts[3] = strconv.Itoa(build)
	newVersion := strings.Join(parts, ".")

	if err := os.WriteFile(path, []byte(newVersion+"\n"), 0644); err != nil {
		return "", err
	}

	return newVersion, nil
}
