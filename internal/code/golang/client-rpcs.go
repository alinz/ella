package golang

import (
	"io"
	"sort"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
	"ella.to/pkg/strcase"
)

func generateClientRpcs(out io.Writer, services []*ast.Service, isMessage func(string) bool) error {
	tmpl, err := code.LoadTemplate(files, "client-rpcs")
	if err != nil {
		return err
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].Name.String() < services[j].Name.String()
	})

	templServices := code.Mapper(services, func(service *ast.Service) any {
		serviceName := service.Name.String()

		methods := code.Filter(service.Methods, func(method *ast.Method) bool {
			return method.Type.Val == "rpc"
		})

		sort.Slice(methods, func(i, j int) bool {
			return methods[i].Name.String() < methods[j].Name.String()
		})

		serviceTopicName := "Rpc" + strcase.ToPascal(serviceName) + "Topic"
		serviceTopic := "rpc." + strcase.ToPascal(serviceName) + ".*"

		return struct {
			TopicName string
			Topic     string
			Name      string
			Methods   any
		}{
			TopicName: serviceTopicName,
			Topic:     serviceTopic,
			Name:      strcase.ToPascal(serviceName),
			Methods: code.Mapper(methods, func(method *ast.Method) any {
				methodTopicName := "Rpc" + strcase.ToPascal(serviceName) + strcase.ToPascal(method.Name.String()) + "Topic"
				methodTopic := "rpc." + strcase.ToPascal(serviceName) + "." + strcase.ToPascal(method.Name.String())

				argsWithTypes := code.Reduce(method.Args, func(acc *strings.Builder, arg *ast.Arg, i int) *strings.Builder {
					acc.WriteString(strcase.ToCamel(arg.Name.String()))
					acc.WriteString(" ")
					if isMessage(arg.Type.String()) {
						acc.WriteString("*")
					}
					acc.WriteString(arg.Type.String())
					if i != len(method.Args)-1 {
						acc.WriteString(", ")
					}
					return acc
				}, &strings.Builder{})

				argsDefinitions := code.Reduce(method.Args, func(acc *strings.Builder, arg *ast.Arg, i int) *strings.Builder {
					acc.WriteString(strcase.ToPascal(arg.Name.String()))
					acc.WriteString(" ")
					if isMessage(arg.Type.String()) {
						acc.WriteString("*")
					}
					acc.WriteString(arg.Type.String())
					acc.WriteString("\n")
					return acc
				}, &strings.Builder{})

				argsAssign := code.Reduce(method.Args, func(acc *strings.Builder, arg *ast.Arg, i int) *strings.Builder {
					acc.WriteString(strcase.ToPascal(arg.Name.String()))
					acc.WriteString(": ")
					acc.WriteString(strcase.ToCamel(arg.Name.String()))
					acc.WriteString(",\n")
					return acc
				}, &strings.Builder{})

				returnsWithTypes := code.Reduce(method.Returns, func(acc *strings.Builder, arg *ast.Return, i int) *strings.Builder {
					acc.WriteString(strcase.ToCamel(arg.Name.String()))
					acc.WriteString(" ")
					if isMessage(arg.Type.String()) {
						acc.WriteString("*")
					}
					acc.WriteString(arg.Type.String())
					if i != len(method.Args)-1 {
						acc.WriteString(", ")
					}
					return acc
				}, &strings.Builder{})

				returnsDefinitions := code.Reduce(method.Returns, func(acc *strings.Builder, arg *ast.Return, i int) *strings.Builder {
					acc.WriteString(strcase.ToPascal(arg.Name.String()))
					acc.WriteString(" ")
					if isMessage(arg.Type.String()) {
						acc.WriteString("*")
					}
					acc.WriteString(arg.Type.String())
					acc.WriteString("\n")
					return acc
				}, &strings.Builder{})

				Returns := code.Reduce(method.Returns, func(acc *strings.Builder, arg *ast.Return, i int) *strings.Builder {
					acc.WriteString("out.")
					acc.WriteString(strcase.ToPascal(arg.Name.String()))
					if i != len(method.Args)-1 {
						acc.WriteString(", ")
					}
					return acc
				}, &strings.Builder{})

				return struct {
					Name            string
					TopicName       string
					Topic           string
					ArgsWithTypes   string
					ArgsDefinitions string
					ArgsAssign      string

					ReturnsWithTypes   string
					ReturnsDefinitions string
					Returns            string
				}{
					Name:            strcase.ToPascal(method.Name.String()),
					TopicName:       methodTopicName,
					Topic:           methodTopic,
					ArgsWithTypes:   argsWithTypes.String(),
					ArgsDefinitions: argsDefinitions.String(),
					ArgsAssign:      argsAssign.String(),

					ReturnsWithTypes:   returnsWithTypes.String(),
					ReturnsDefinitions: returnsDefinitions.String(),
					Returns:            Returns.String(),
				}
			}),
		}
	})

	return tmpl.Execute(out, templServices)
}
