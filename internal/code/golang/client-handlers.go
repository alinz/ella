package golang

import (
	"io"
	"sort"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
	"ella.to/pkg/strcase"
)

func generateClientHandlers(out io.Writer, services []*ast.Service, isMessage func(string) bool) error {
	tmpl, err := code.LoadTemplate(files, "client-handlers")
	if err != nil {
		return err
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].Name.String() < services[j].Name.String()
	})

	templServices := code.Mapper(services, func(service *ast.Service) any {
		serviceName := service.Name.String()

		methods := code.Filter(service.Methods, func(method *ast.Method) bool {
			return method.Type.Val == "http"
		})

		sort.Slice(methods, func(i, j int) bool {
			return methods[i].Name.String() < methods[j].Name.String()
		})

		servicePathName := strcase.ToPascal(serviceName) + "HandlerServicePrefixPath"

		return struct {
			Name     string
			NameImpl string
			PathName string
			Methods  any
		}{
			Name:     strcase.ToPascal(serviceName) + "HandlerService",
			NameImpl: strcase.ToCamel(serviceName) + "HandlerServiceClient",
			PathName: servicePathName,
			Methods: code.Mapper(methods, func(method *ast.Method) any {

				sbArgsInput := code.Reduce(method.Args, func(acc *strings.Builder, arg *ast.Arg, i int) *strings.Builder {
					name := arg.Name.String()
					typ := parseType(arg.Type)
					pointer := isMessage(arg.Type.String())
					acc.WriteString(writeVariableWithType(name, typ, pointer))
					if i < len(method.Args)-1 {
						acc.WriteString(", ")
					}
					return acc
				}, &strings.Builder{})

				sbArgsWithTypes := code.Reduce(method.Args, func(acc *strings.Builder, arg *ast.Arg, i int) *strings.Builder {
					name := arg.Name.String()
					typ := parseType(arg.Type)
					pointer := isMessage(arg.Type.String())
					acc.WriteString(writeVariableWithType(strcase.ToPascal(name), typ, pointer))
					if i < len(method.Args)-1 {
						acc.WriteString("\n")
					}
					return acc
				}, &strings.Builder{})

				sbArgs := code.Reduce(method.Args, func(acc *strings.Builder, arg *ast.Arg, i int) *strings.Builder {
					name := arg.Name.String()
					acc.WriteString(strcase.ToPascal(name))
					acc.WriteString(": ")
					acc.WriteString(name)
					acc.WriteString(", \n")
					return acc
				}, &strings.Builder{})

				sbReturnsWithTypes := code.Reduce(method.Returns, func(acc *strings.Builder, arg *ast.Return, i int) *strings.Builder {
					name := arg.Name.String()
					typ := parseType(arg.Type)
					pointer := isMessage(arg.Type.String())
					acc.WriteString(writeVariableWithType(name, typ, pointer))
					if i < len(method.Returns)-1 {
						acc.WriteString("\n")
					}
					return acc
				}, &strings.Builder{})

				sbReturns := code.Reduce(method.Returns, func(acc *strings.Builder, arg *ast.Return, i int) *strings.Builder {
					name := arg.Name.String()
					acc.WriteString("out.")
					acc.WriteString(name)
					if i < len(method.Returns)-1 {
						acc.WriteString(", ")
					}
					return acc
				}, &strings.Builder{})

				hasArgs := "true"
				if len(method.Args) == 0 {
					hasArgs = "false"
				}

				return struct {
					Name             string
					Http             string
					PathName         string
					HasArgs          string
					ArgsInput        string
					ArgsWithTypes    string
					Args             string
					ReturnsWithTypes string
					Returns          string
				}{
					Name:             strcase.ToPascal(method.Name.String()),
					Http:             getHttpMethodCall(method),
					PathName:         strcase.ToPascal(serviceName) + strcase.ToPascal(method.Name.String()) + "MethodPrefixPath",
					HasArgs:          hasArgs,
					ArgsInput:        sbArgsInput.String(),
					ArgsWithTypes:    sbArgsWithTypes.String(),
					Args:             sbArgs.String(),
					ReturnsWithTypes: sbReturnsWithTypes.String(),
					Returns:          sbReturns.String(),
				}
			}),
		}
	})

	return tmpl.Execute(out, templServices)
}
