package main

import (
	"fmt"
	"os"
	"path/filepath"

	"ella.to/internal/parser"
)

const usage = `Usage: ella [command]

Commands:
  - fmt Format a file in place
        ella fmt <glob path>

  - gen Generate code from a folder to a file and
        currently supports .go and .ts 
        ella gen <glob path> <output path to file>

example:
  ella fmt ./path/to/*.ella
  ella gen ./path/to/*.ella ./path/to/output.go
  ella gen ./path/to/*.ella ./path/to/output.ts
`

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(usage)
		os.Exit(0)
	}

	var err error

	switch os.Args[1] {
	case "fmt":
		if len(os.Args) != 3 {
			fmt.Printf(usage)
			os.Exit(0)
		}
		err = format(os.Args[2])
	case "gen":
		if len(os.Args) != 4 {
			fmt.Printf(usage)
			os.Exit(0)
		}
		err = gen(os.Args[2], os.Args[3])
	default:
		fmt.Printf(usage)
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func format(path string) error {
	matches, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	for _, match := range matches {
		content, err := os.ReadFile(match)
		if err != nil {
			return err
		}

		prog, err := parser.ParseProgram(
			parser.New(string(content)),
		)
		if err != nil {
			return err
		}

		err = os.WriteFile(match, []byte(prog.String()), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func gen(globIn, outputFile string) error {
	return nil
}
