package golang

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"

	"ella.to/internal/ast"
	"ella.to/internal/code"
)

//go:embed templates/*.tmpl
var files embed.FS

type Golang struct {
	PkgName      string
	Constants    Constants
	Enums        Enums
	Bases        Bases
	Messages     Messages
	HttpServices HttpServices
	RpcServices  RpcServices
}

func (g *Golang) Parse(prog *ast.Program) error {
	return code.RunParsers(
		prog,
		g.Constants.Parse,
		g.Enums.Parse,
		g.Bases.Parse,
		g.Messages.Parse,
		g.HttpServices.Parse,
		g.RpcServices.Parse,
	)
}

func New(pkg string) code.Generator {
	return code.GeneratorFunc(func(outFilename string, prog *ast.Program) error {
		golang := Golang{
			PkgName: pkg,
		}

		if err := golang.Parse(prog); err != nil {
			return err
		}

		tmpl, err := code.LoadTemplate(files, "templates", "golang")
		if err != nil {
			return err
		}

		out, err := os.Create(outFilename)
		if err != nil {
			return err
		}
		defer out.Close()

		if err := tmpl.Execute(out, golang); err != nil {
			return err
		}

		var errBuffer bytes.Buffer
		formatCmd := exec.Command("go", "fmt", outFilename)
		formatCmd.Stderr = &errBuffer
		if err = formatCmd.Run(); err != nil {
			return fmt.Errorf("%s: %s", err, errBuffer.String())
		}

		return nil
	})
}

func castString(value any, defaultValue string) string {
	return castValue[string](value, defaultValue)
}

func castInt64(value any, defaultValue int64) int64 {
	return castValue[int64](value, defaultValue)
}

func castBool(value any, defaultValue bool) bool {
	return castValue[bool](value, defaultValue)
}

func castValue[T any](value any, defaultValue T) T {
	result, ok := value.(T)
	if ok {
		return result
	}

	return defaultValue
}

func isArrayOf[T ast.Type](typ ast.Type) bool {
	arr, ok := typ.(*ast.Array)
	if !ok {
		return false
	}

	_, ok = arr.Type.(T)
	return ok
}

func parseValueType(value ast.Value) string {
	switch v := value.(type) {
	case *ast.ValueByteSize:
		return "int64"
	case *ast.ValueDuration:
		return "int64"
	case *ast.ValueInt:
		return fmt.Sprintf("int%d", v.Size)
	case *ast.ValueUint:
		return fmt.Sprintf("uint%d", v.Size)
	case *ast.ValueFloat:
		return fmt.Sprintf("float%d", v.Size)
	case *ast.ValueString:
		return "string"
	case *ast.ValueBool:
		return "bool"
	case *ast.ValueNull:
		return "any"
	case *ast.ValueVariable:
		return "any" // TODO: find a way to get the type in recersive mode, be aware of cycles
	default:
		panic(fmt.Errorf("unknown type for value: %T", value))
	}
}
