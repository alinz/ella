package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/alinz/rpc.go/schema/parser"
	"github.com/alinz/rpc.go/schema/validator"
	"github.com/alinz/rpc.go/templates/golang"
)

func genCmd() *cli.Command {
	var outDir string
	var schemaFile string

	return &cli.Command{
		Name:  "gen",
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
				Usage:       "target's input schema file `./example/schema.rpc`",
				Required:    true,
				Destination: &schemaFile,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			err = os.Mkdir(outDir, 0755)
			if err != nil && !errors.Is(err, os.ErrExist) {
				return err
			}

			in, err := os.Open(schemaFile)
			if err != nil {
				return err
			}
			defer in.Close()
			inData, err := io.ReadAll(in)
			if err != nil {
				return err
			}

			program, err := parser.New(string(inData)).Parse()
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
			genCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
