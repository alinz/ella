package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"ella.to/schema/parser"
	"ella.to/schema/validator"
	"ella.to/templates/golang"
)

func fmtCmd() *cli.Command {
	var schemaFile string

	return &cli.Command{
		Name:  "fmt",
		Usage: "format rpc schema file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "target's input schema file `./example/schema.rpc`",
				Required:    true,
				Destination: &schemaFile,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			in, err := os.Open(schemaFile)
			if err != nil {
				return err
			}
			inData, err := io.ReadAll(in)
			if err != nil {
				in.Close()
				return err
			}
			in.Close()

			program, err := parser.New(string(inData)).Parse()
			if err != nil {
				return err
			}

			err = validator.Validate(program)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "%s", program.TokenLiteral())

			return nil
		},
	}
}

func rpcCmd() *cli.Command {
	var outDir string
	var schemaDir string

	return &cli.Command{
		Name:  "rpc",
		Usage: "generate rpc client and server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "target's output directory `./example`",
				Required:    true,
				Destination: &outDir,
			},
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "target's input schema folder `./example/schema/`",
				Required:    true,
				Destination: &schemaDir,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			err = os.Mkdir(outDir, 0755)
			if err != nil && !errors.Is(err, os.ErrExist) {
				return err
			}

			var inData bytes.Buffer

			{
				files, err := os.ReadDir(schemaDir)
				if err != nil {
					return err
				}

				for _, file := range files {
					if file.IsDir() {
						continue
					}

					if filepath.Ext(file.Name()) != ".ella" {
						continue
					}

					err = func(buffer *bytes.Buffer, filename string) error {
						in, err := os.Open(filename)
						if err != nil {
							return err
						}
						defer in.Close()

						_, err = io.Copy(buffer, in)
						if err != nil {
							return err
						}

						buffer.WriteString("\n")
						return nil
					}(&inData, filepath.Join(schemaDir, file.Name()))
					if err != nil {
						return err
					}
				}
			}

			program, err := parser.New(inData.String()).Parse()
			if err != nil {
				return err
			}

			err = validator.Validate(program)
			if err != nil {
				return err
			}

			outFile := filepath.Join(outDir, "rpc.gen.go")
			out, err := os.Create(outFile)
			if err != nil {
				return err
			}
			defer out.Close()

			err = golang.Generate(out, "rpc", program)
			if err != nil {
				return err
			}

			err = exec.Command("go", "fmt", outFile).Run()
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			rpcCmd(),
			fmtCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
