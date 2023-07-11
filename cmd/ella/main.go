package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
	"ella.to/internal/code/golang"
	"ella.to/internal/code/typescript"
	"ella.to/internal/parser"
)

const usage = `Usage: ella [command]

Commands:
  - fmt Format one or many files in place using glob pattern
        ella fmt <glob path>

  - gen Generate code from a folder to a file and
        currently supports .go and .ts 
        ella gen <pkg> <search glob path> <output path to file>

  - ver Print the version of ella

example:
  ella fmt ./path/to/*.ella
  ella gen rpc ./path/to/*.ella ./path/to/output.go
  ella gen rpc ./path/to/*.ella ./path/to/output.ts
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(0)
	}

	var err error

	switch os.Args[1] {
	case "fmt":
		if len(os.Args) != 3 {
			fmt.Print(usage)
			os.Exit(0)
		}
		err = format(os.Args[2])
	case "gen":
		if len(os.Args) != 5 {
			fmt.Print(usage)
			os.Exit(0)
		}
		err = gen(os.Args[2], os.Args[3], os.Args[4])
	default:
		fmt.Print(usage)
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func format(path string) error {
	filenames, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	for _, filename := range filenames {
		prog, err := parse(filename)
		if err != nil {
			return err
		}

		err = os.WriteFile(filename, []byte(prog.String()), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func gen(pkg, search, out string) error {
	var code code.Generator

	ext := filepath.Ext(out)

	switch ext {
	case ".go":
		code = golang.New()
	case ".ts":
		code = typescript.New()
	default:
		return fmt.Errorf("unknown extension %s", out)
	}

	filenames, err := filepath.Glob(search)
	if err != nil {
		return err
	}

	content, err := combine(filenames...)
	if err != nil {
		return err
	}

	prog, err := parser.ParseProgram(parser.New(content))
	if err != nil {
		return err
	}

	fileout, err := os.Create(out)
	if err != nil {
		return err
	}

	err = code.Generate(fileout, pkg, prog)
	if err != nil {
		return err
	}

	if ext == ".go" {
		err = exec.Command("go", "fmt", out).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func parse(filename string) (*ast.Program, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return parser.ParseProgram(
		parser.New(string(content)),
	)
}

// combine is a helper function to concatenate multiple files into a single string
// and returns an error if any of the files have an invalid extension or file cannot
// be read.
func combine(filenames ...string) (string, error) {
	var sb strings.Builder

	for i, filename := range filenames {
		if !strings.HasSuffix(filename, ".ella") {
			return "", fmt.Errorf("invalid file extension %s", filename)
		}

		content, err := os.ReadFile(filename)
		if err != nil {
			return "", err
		}

		if i != 0 {
			sb.WriteString("\n")
		}
		sb.Write(content)
	}

	return sb.String(), nil
}
