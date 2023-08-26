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
	Messages     Messages
	HttpServices HttpServices
	RpcServices  RpcServices
}

func (g *Golang) Parse(prog *ast.Program) error {
	return code.RunParsers(
		prog,
		g.Constants.Parse,
		g.Enums.Parse,
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

func createIsMessageTypeFunc(messages []*ast.Message) func(value string) bool {
	messagesMap := make(map[string]struct{})
	for _, message := range messages {
		messagesMap[message.Name.String()] = struct{}{}
	}

	return func(value string) bool {
		_, ok := messagesMap[value]
		return ok
	}
}
