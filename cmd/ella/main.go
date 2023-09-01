package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
	"ella.to/internal/code/golang"
	"ella.to/internal/code/typescript"
	"ella.to/internal/parser"
	"ella.to/internal/validator"
)

const Version = "0.0.2"

const usage = `
┌─┐┬  ┬  ┌─┐
├┤ │  │  ├─┤
└─┘┴─┘┴─┘┴ ┴ v` + Version + `

Usage: ella [command]

Commands:
  - fmt Format one or many files in place using glob pattern
        ella fmt <glob path>

  - gen Generate code from a folder to a file and currently
        supports .go and .ts extensions
        ella gen <pkg> <output path to file> <search glob paths...>

  - ver Print the version of ella

example:
  ella fmt ./path/to/*.ella
  ella gen rpc ./path/to/output.go ./path/to/*.ella
  ella gen rpc ./path/to/output.ts ./path/to/*.ella ./path/to/other/*.ella
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
		if len(os.Args) < 5 {
			fmt.Print(usage)
			os.Exit(0)
		}
		err = gen(os.Args[2], os.Args[3], os.Args[4:]...)
	case "ver":
		fmt.Println(Version)
	default:
		fmt.Print(usage)
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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

func gen(pkg, out string, searchPaths ...string) (err error) {
	var code code.Generator

	defer func() {
		if err != nil {
			os.Remove(out)
		}
	}()

	filenames, err := mergeAllFiles(searchPaths...)
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

	err = validator.Validate(prog)
	if err != nil {
		return err
	}

	ext := filepath.Ext(out)
	switch ext {
	case ".go":
		code = golang.New(pkg)
	case ".ts":
		code = typescript.New()
	default:
		return fmt.Errorf("unknown extension %s", out)
	}

	if err = code.Generate(out, prog); err != nil {
		return err
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

func mergeAllFiles(paths ...string) ([]string, error) {
	filenamesMap := make(map[string]struct{})

	for _, path := range paths {
		filenames, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		for _, filename := range filenames {
			filenamesMap[filename] = struct{}{}
		}
	}

	filenames := make([]string, 0, len(filenamesMap))
	for filename := range filenamesMap {
		filenames = append(filenames, filename)
	}

	return filenames, nil
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
