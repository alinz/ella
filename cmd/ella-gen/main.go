package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"ella.to/schema/ast"
	"ella.to/schema/parser"
	"ella.to/schema/validator"
	"ella.to/templates/golang"
	"ella.to/transform"
	"ella.to/transform/restclient"
	"ella.to/transform/typescript"
)

var Version = "main"
var GitCommit = "development"

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Printf("ella-gen version %s.%s\n", cCtx.App.Version, GitCommit)
	}

	app := &cli.App{
		Name:    "ella-gen",
		Usage:   "generate common code that you don't want to write yourself",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "target's input schema folder `./example/schema`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "target's output file `./example/rpc.gen.go`",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "pkg",
				Aliases: []string{"p"},
				Usage:   "target's output package `example` for golang",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			output := ctx.String("output")
			input := ctx.String("input")
			pkg := ctx.String("pkg")

			ext, err := checkFileExtension(output, []string{".go", ".ts", ".http"})
			if err != nil {
				return err
			}

			outputDir := filepath.Dir(output)

			err = os.Mkdir(outputDir, 0755)
			if err != nil && !errors.Is(err, os.ErrExist) {
				return err
			}

			content, err := combineFiles(input, ".ella")
			if err != nil {
				return err
			}

			program, err := parser.New(content).Parse()
			if err != nil {
				return err
			}

			err = validator.Validate(program)
			if err != nil {
				return err
			}

			switch ext {
			case ".go":
				err = golangGen(output, program, pkg)
			case ".ts":
				err = typescriptGen(output, program)
			case ".http":
				err = restClientGen(output, program)
			}

			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func golangGen(outputFile string, program *ast.Program, pkg string) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	err = golang.Generate(out, pkg, program)
	if err != nil {
		return err
	}

	err = exec.Command("go", "fmt", outputFile).Run()
	if err != nil {
		return err
	}

	return nil
}

func typescriptGen(output string, program *ast.Program) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	transform.Run(
		out,
		typescript.Signature(),
		typescript.Constants(ast.GetSlice[*ast.Constant](program)),
		typescript.Enums(ast.GetSlice[*ast.Enum](program)),
		typescript.Messages(ast.GetSlice[*ast.Message](program)),
		typescript.Services(ast.GetSlice[*ast.Service](program)),
		typescript.HelperFunc(),
	)

	return nil
}

func restClientGen(output string, program *ast.Program) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	transform.Run(
		out,
		restclient.Constants(program),
		restclient.Services(program),
	)

	return nil
}
