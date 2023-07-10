package golang

import (
	"io"
	"sort"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
)

func writeVariableWithType(name string, typ string, isPointer bool) string {
	var sb strings.Builder

	sb.WriteString(name)
	sb.WriteString(" ")
	if isPointer {
		sb.WriteString("*")
	}
	sb.WriteString(typ)

	return sb.String()
}

func generateServices(out io.Writer, services []*ast.Service, isMessage func(string) bool) error {
	tmpl, err := code.LoadTemplate(files, "services")
	if err != nil {
		return err
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].Name.String() < services[j].Name.String()
	})

	templServices := code.Mapper(services, func(service *ast.Service) any {
		return struct {
			Name    string
			Methods []any
		}{
			Name: service.Name.String(),
			Methods: code.Mapper(service.Methods, func(method *ast.Method) any {
				returns := code.Mapper(method.Returns, func(arg *ast.Return) string {
					name := arg.Name.String()
					typ := parseType(arg.Type)
					pointer := isMessage(arg.Type.String())
					return writeVariableWithType(name, typ, pointer)
				})

				// every method has an error return at the last position
				returns = append(returns, "err error")

				args := code.Mapper(method.Args, func(arg *ast.Arg) string {
					name := arg.Name.String()
					typ := parseType(arg.Type)
					pointer := isMessage(arg.Type.String())
					return writeVariableWithType(name, typ, pointer)
				})

				// every method's args has a context as the first argument
				args = append([]string{"ctx context.Context"}, args...)

				return struct {
					Type    string
					Name    string
					Args    string
					Returns string
				}{
					Type:    method.Type.Val,
					Name:    method.Name.String(),
					Args:    strings.Join(args, ", "),
					Returns: strings.Join(returns, ", "),
				}
			}),
		}
	})

	return tmpl.Execute(out, templServices)
}
