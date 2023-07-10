package golang

import (
	"embed"
	"fmt"
	"io"

	"ella.to/internal/ast"
	"ella.to/internal/code"
)

//go:embed tmpl/*.tmpl
var files embed.FS

func New() code.Generator {
	return code.GeneratorFunc(func(w io.Writer, pkg string, prog *ast.Program) (err error) {
		constants := code.GetConstants(prog)
		enums := code.GetEnums(prog)
		messages := code.GetMessages(prog)
		services := code.GetServices(prog)

		isMessage := func() func(name string) bool {
			messagesSet := make(map[string]struct{})
			for _, message := range messages {
				messagesSet[message.Name.String()] = struct{}{}
			}

			return func(name string) bool {
				_, ok := messagesSet[name]
				return ok
			}
		}()

		err = generateHeader(w, pkg)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\n\n")

		err = generateConstants(w, pkg, constants)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\n\n")

		err = generateEnums(w, enums)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\n")

		err = generateMessages(w, messages, isMessage)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\n")

		err = generateServices(w, services, isMessage)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\n")

		err = generateServerHandlers(w, services, isMessage)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\n")

		err = generateClientHandlers(w, services, isMessage)
		if err != nil {
			return err
		}

		fmt.Fprint(w, "\n")

		err = generateHelper(w)
		if err != nil {
			return err
		}

		return nil
	})
}
