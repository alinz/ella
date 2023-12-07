package golang

import (
	"fmt"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/ast/astutil"
	"ella.to/pkg/sliceutil"
	"ella.to/pkg/strcase"
)

type MethodArg struct {
	Name string
	Type string
}

type MethodArgs []MethodArg

func (m MethodArgs) Definitions() string {
	return strings.Join(sliceutil.Mapper(m, func(arg MethodArg) string {
		return fmt.Sprintf(", %s %s", arg.Name, arg.Type)
	}), "")
}

type MethodReturn struct {
	Name   string
	Type   string
	Stream bool
}

type MethodReturns []MethodReturn

func (m MethodReturns) Definitions() string {
	m = append(m, MethodReturn{
		Name: "err",
		Type: "error",
	})

	return strings.Join(sliceutil.Mapper(m, func(arg MethodReturn) string {
		return fmt.Sprintf("%s %s", arg.Name, arg.Type)
	}), ", ")
}

type Method struct {
	Name    string
	Service string
	Options astutil.MethodOptions
	Args    MethodArgs
	Returns MethodReturns
}

func (m Method) GetReturnStreamName() (string, error) {
	for _, ret := range m.Returns {
		if ret.Stream {
			return strcase.ToCamel(ret.Name), nil
		}
	}

	return "", fmt.Errorf("method has no stream return")
}

func (m Method) HasArgs() string {
	if len(m.Args) > 0 {
		return "true"
	}

	return "false"
}

func (m Method) HasReturns() bool {
	return len(m.Returns) > 0
}

func (m Method) TopicName() string {
	return fmt.Sprintf("TopicRpc%s%sMethod", strcase.ToPascal(m.Service), strcase.ToPascal(m.Name))
}

func (m Method) TopicValue() string {
	return fmt.Sprintf("ella.rpc.%s.%s", strcase.ToSnake(m.Service), strcase.ToSnake(m.Name))
}

func (m Method) PathName() string {
	return fmt.Sprintf("PathHttp%s%sMethod", strcase.ToPascal(m.Service), strcase.ToPascal(m.Name))
}

func (m Method) PathValue() string {
	return fmt.Sprintf("/ella/http/%s/%s", strcase.ToPascal(m.Service), strcase.ToPascal(m.Name))
}

func (m Method) ArgsNames(prefix string) string {
	return strings.Join(sliceutil.Mapper(sliceutil.Filter(m.Args, func(arg MethodArg) bool {
		return arg.Type != "func() (string, io.Reader, error)"
	}), func(arg MethodArg) string {
		return prefix + strcase.ToPascal(arg.Name) + ","
	}), "\n")
}

func (m Method) ReturnStreamType() string {
	for _, ret := range m.Returns {
		if ret.Stream {
			return strings.Replace(ret.Type, "<-chan ", "", 1)
		}
	}

	return ""
}

func (m Method) ArgsStructDefinitions(pointer bool) string {
	return strings.Join(sliceutil.Mapper(sliceutil.Filter(m.Args, func(arg MethodArg) bool {
		return arg.Type != "func() (string, io.Reader, error)"
	}), func(arg MethodArg) string {
		return fmt.Sprintf("%s %s `json:\"%s\"`", strcase.ToPascal(arg.Name), arg.Type, strcase.ToSnake(arg.Name))
	}), "\n")
}

func (m Method) ArgsNamesValues() string {
	return strings.Join(sliceutil.Mapper(sliceutil.Filter(m.Args, func(arg MethodArg) bool {
		return arg.Type != "func() (string, io.Reader, error)"
	}), func(arg MethodArg) string {
		return strcase.ToPascal(arg.Name) + ":" + arg.Name + ","
	}), "\n")
}

func (m Method) ReturnsNames(prefix string) string {
	return strings.Join(sliceutil.Mapper(m.Returns, func(ret MethodReturn) string {
		return prefix + strcase.ToPascal(ret.Name) + ", "
	}), "")
}

func (m Method) ReturnsStructDefinitions() string {
	return strings.Join(sliceutil.Mapper(m.Returns, func(ret MethodReturn) string {
		return fmt.Sprintf("%s %s `json:\"%s\"`", strcase.ToPascal(ret.Name), ret.Type, strcase.ToSnake(ret.Name))
	}), "\n")
}

func (m Method) IsStream() bool {
	for _, ret := range m.Returns {
		if ret.Stream {
			return true
		}
	}

	return false
}

func (m Method) IsBinary() bool {
	for _, ret := range m.Returns {
		if ret.Type == "io.Reader" {
			return true
		}
	}

	return false
}

func (m Method) IsFileUpload() bool {
	var result bool

	for _, arg := range m.Args {
		if arg.Type == "func() (string, io.Reader, error)" {
			result = true
			break
		}
	}

	return result
}

type Methods []Method

type HttpService struct {
	Name    string
	Methods Methods
}

func (s HttpService) NameImpl() string {
	return "http" + strcase.ToPascal(s.Name) + "Server"
}

func (s HttpService) PathName() string {
	return fmt.Sprintf("PathHttp%sPrefix", strcase.ToPascal(s.Name))
}

func (s HttpService) PathValue() string {
	return fmt.Sprintf("/ella/http/%s/", strcase.ToPascal(s.Name))
}

type HttpServices []HttpService

func (s *HttpServices) Parse(prog *ast.Program) error {
	isModelType := astutil.CreateIsModelTypeFunc(astutil.GetModels(prog))

	*s = sliceutil.Mapper(astutil.GetServices(prog), func(service *ast.Service) HttpService {
		methods := sliceutil.Filter(service.Methods, func(method *ast.Method) bool {
			return method.Type == ast.MethodHTTP
		})

		return HttpService{
			Name: service.Name.String(),
			Methods: sliceutil.Mapper(methods, func(method *ast.Method) Method {
				return Method{
					Name:    method.Name.String(),
					Service: service.Name.String(),
					Options: astutil.ParseMethodOptions(method.Options),
					Args: sliceutil.Mapper(method.Args, func(arg *ast.Arg) MethodArg {
						var typ string

						if _, ok := arg.Type.(*ast.File); ok {
							typ = "func() (string, io.Reader, error)"
						} else {
							typ = parseType(arg.Type, isModelType)
						}

						return MethodArg{
							Name: arg.Name.String(),
							Type: typ,
						}
					}),
					Returns: sliceutil.Mapper(method.Returns, func(ret *ast.Return) MethodReturn {
						typ := parseType(ret.Type, isModelType)
						if ret.Stream && isArrayOf[*ast.Byte](ret.Type) {
							typ = "io.Reader"
						} else if ret.Stream {
							typ = "<-chan " + typ
						}

						return MethodReturn{
							Name:   ret.Name.String(),
							Type:   typ,
							Stream: ret.Stream,
						}
					}),
				}
			}),
		}
	})

	// we want to make sure that we don't generate services without methods
	*s = sliceutil.Filter(*s, func(service HttpService) bool {
		return len(service.Methods) != 0
	})

	return nil
}

type RpcService struct {
	Name    string
	Methods Methods
}

func (s RpcService) TopicName() string {
	return fmt.Sprintf("TopicRpc%s", strcase.ToPascal(s.Name))
}

func (s RpcService) TopicValue() string {
	return fmt.Sprintf("ella.rpc.%s.*", strcase.ToSnake(s.Name))
}

type RpcServices []RpcService

func (s *RpcServices) Parse(prog *ast.Program) error {
	*s = sliceutil.Mapper(astutil.GetServices(prog), func(service *ast.Service) RpcService {
		isModelType := astutil.CreateIsModelTypeFunc(astutil.GetModels(prog))

		methods := sliceutil.Filter(service.Methods, func(method *ast.Method) bool {
			return method.Type == ast.MethodRPC
		})

		return RpcService{
			Name: service.Name.String(),
			Methods: sliceutil.Mapper(methods, func(method *ast.Method) Method {
				return Method{
					Name:    method.Name.String(),
					Service: service.Name.String(),
					Args: sliceutil.Mapper(method.Args, func(arg *ast.Arg) MethodArg {
						return MethodArg{
							Name: arg.Name.String(),
							Type: parseType(arg.Type, isModelType),
						}
					}),
					Returns: sliceutil.Mapper(method.Returns, func(ret *ast.Return) MethodReturn {
						return MethodReturn{
							Name: ret.Name.String(),
							Type: parseType(ret.Type, isModelType),
						}
					}),
				}
			}),
		}
	})

	// we want to make sure that we don't generate services without methods
	*s = sliceutil.Filter(*s, func(service RpcService) bool {
		return len(service.Methods) != 0
	})

	return nil
}
