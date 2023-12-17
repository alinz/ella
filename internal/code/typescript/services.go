package typescript

import (
	"fmt"

	"ella.to/internal/ast"
	"ella.to/internal/ast/astutil"
	"ella.to/pkg/sliceutil"
	"ella.to/pkg/strcase"
)

type Arg struct {
	Name string
	Type string
}

type Return struct {
	Name string
	Type string
}

type Method struct {
	Name        string
	ServiceName string
	Options     astutil.MethodOptions
	Type        string // normal, binary, stream, fileupload
	Args        []Arg
	Returns     []Return
}

func (m Method) PathValue() string {
	return fmt.Sprintf("/ella/http/%s/%s", strcase.ToPascal(m.ServiceName), strcase.ToPascal(m.Name))
}

func (m Method) ArgsName() string {
	return fmt.Sprintf("Service%s%sArgs", m.ServiceName, strcase.ToPascal(m.Name))
}

func (m Method) ReturnsName() string {
	switch m.Type {
	case "binary":
		return "Blob"
	case "stream":
		return "Subscription<" + m.Returns[0].Type + ">"
	default:
		return fmt.Sprintf("Service%s%sReturns", m.ServiceName, strcase.ToPascal(m.Name))
	}
}

func (m Method) HasReturn() bool {
	return !(m.IsBinaryStream() || m.IsStream())
}

func (m Method) IsStream() bool {
	return m.Type == "stream"
}

func (m Method) IsBinaryStream() bool {
	return m.Type == "binary"
}

func (m Method) IsFileUpload() bool {
	return m.Type == "fileupload"
}

func (m Method) NeedReturnInterface() bool {
	return !(m.IsStream() || m.IsBinaryStream())
}

type HttpService struct {
	Name    string
	Methods []Method
}

type HttpServices []HttpService

func (s *HttpServices) Parse(prog *ast.Program) error {
	*s = sliceutil.Mapper(astutil.GetServices(prog), func(service *ast.Service) HttpService {
		return HttpService{
			Name: service.Name.String(),
			Methods: sliceutil.Mapper(sliceutil.Filter(service.Methods, func(method *ast.Method) bool {
				return method.Type == ast.MethodHTTP
			}), func(method *ast.Method) Method {

				m := Method{}

				m.ServiceName = service.Name.String()
				m.Name = strcase.ToCamel(method.Name.String())
				m.Options = astutil.ParseMethodOptions(method.Options)

				m.Args = sliceutil.Mapper(sliceutil.Filter(
					method.Args,
					func(arg *ast.Arg) bool {
						if m.Type != "fileupload" && arg.Type.String() == "file" {
							m.Type = "fileupload"
							return false
						}

						return arg.Type.String() != "file"
					},
				), func(arg *ast.Arg) Arg {
					return Arg{
						Name: strcase.ToSnake(arg.Name.String()),
						Type: parseType(arg.Type),
					}
				})
				m.Returns = sliceutil.Mapper(method.Returns, func(ret *ast.Return) Return {
					typ := parseType(ret.Type)
					if ret.Stream && typ == "byte[]" {
						m.Type = "binary"
					} else if ret.Stream {
						m.Type = "stream"
					} else if m.Type != "fileupload" {
						m.Type = "normal"
					}

					return Return{
						Name: strcase.ToSnake(ret.Name.String()),
						Type: typ,
					}
				})

				return m
			}),
		}
	})

	return nil
}
