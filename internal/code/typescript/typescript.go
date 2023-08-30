package typescript

import (
	"embed"
	"os"

	"ella.to/internal/ast"
	"ella.to/internal/code"
)

//go:embed templates/*.tmpl
var files embed.FS

type Typescript struct {
	Constants Constants
	Enums     Enums
	Models    Models
	// HttpServices HttpServices
}

func (t *Typescript) Parse(prog *ast.Program) error {
	return code.RunParsers(
		prog,
		t.Constants.Parse,
		t.Enums.Parse,
		t.Models.Parse,
		// t.HttpServices.Parse,
		// t.RpcServices.Parse,
	)
}

func New() code.Generator {
	return code.GeneratorFunc(func(outFilename string, prog *ast.Program) error {
		typescript := Typescript{}

		if err := typescript.Parse(prog); err != nil {
			return err
		}

		tmpl, err := code.LoadTemplate(files, "templates", "typescript")
		if err != nil {
			return err
		}

		out, err := os.Create(outFilename)
		if err != nil {
			return err
		}
		defer out.Close()

		if err := tmpl.Execute(out, typescript); err != nil {
			return err
		}

		return nil
	})
}
