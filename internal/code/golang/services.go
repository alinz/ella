package golang

import (
	"fmt"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/astutil"
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

type MethodOptions struct {
	HttpMethod    string
	MaxUploadSize int64
	RawControl    bool
}

type Method struct {
	Name    string
	Service string
	Options MethodOptions
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
	return fmt.Sprintf("TopicRpc%sService%sMethod", strcase.ToPascal(m.Service), strcase.ToPascal(m.Name))
}

func (m Method) TopicValue() string {
	return fmt.Sprintf("ella.rpc.%s.%s", strcase.ToSnake(m.Service), strcase.ToSnake(m.Name))
}

func (m Method) PathName() string {
	return fmt.Sprintf("PathHttp%sService%sMethod", strcase.ToPascal(m.Service), strcase.ToPascal(m.Name))
}

func (m Method) PathValue() string {
	return fmt.Sprintf("/ella/http/%s/%s", strcase.ToPascal(m.Service), strcase.ToPascal(m.Name))
}

func (m Method) ArgsNames(prefix string) string {
	return strings.Join(sliceutil.Mapper(sliceutil.Filter(m.Args, func(arg MethodArg) bool {
		return arg.Type != "<-chan *fileUpload"
	}), func(arg MethodArg) string {
		return prefix + strcase.ToPascal(arg.Name) + ","
	}), "\n")
}

func (m Method) ArgsStructDefinitions() string {
	return strings.Join(sliceutil.Mapper(sliceutil.Filter(m.Args, func(arg MethodArg) bool {
		return arg.Type != "<-chan *fileUpload"
	}), func(arg MethodArg) string {
		return fmt.Sprintf("%s %s", strcase.ToPascal(arg.Name), arg.Type)
	}), "\n")
}

func (m Method) ArgsNamesValues() string {
	return strings.Join(sliceutil.Mapper(m.Args, func(arg MethodArg) string {
		return strcase.ToPascal(arg.Name) + ":" + arg.Name + ","
	}), "\n")
}

func (m Method) ReturnsNames(prefix string) string {
	return strings.Join(sliceutil.Mapper(m.Returns, func(arg MethodReturn) string {
		return prefix + strcase.ToPascal(arg.Name) + ", "
	}), "")
}

func (m Method) ReturnsStructDefinitions() string {
	return strings.Join(sliceutil.Mapper(m.Returns, func(arg MethodReturn) string {
		return fmt.Sprintf("%s %s", strcase.ToPascal(arg.Name), arg.Type)
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

func (m Method) IsFileUpload() bool {
	var result bool

	for _, arg := range m.Args {
		if arg.Type == "<-chan *fileUpload" {
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
	return "http" + strcase.ToPascal(s.Name) + "ServiceServer"
}

func (s HttpService) PathName() string {
	return fmt.Sprintf("PathHttp%sServicePrefix", strcase.ToPascal(s.Name))
}

func (s HttpService) PathValue() string {
	return fmt.Sprintf("/ella/http/%s/", strcase.ToPascal(s.Name))
}

type HttpServices []HttpService

func (s *HttpServices) Parse(prog *ast.Program) error {
	isMessageType := createIsMessageTypeFunc(astutil.GetMessages(prog))

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
					Options: parseMethodOptions(method),
					Args: sliceutil.Mapper(method.Args, func(arg *ast.Arg) MethodArg {
						var typ string

						if _, ok := arg.Type.(*ast.File); ok {
							typ = "<-chan *fileUpload"
						} else {
							typ = arg.Type.String()
							if isMessageType(arg.Type.String()) {
								typ = fmt.Sprintf("*%s", typ)
							}
						}

						return MethodArg{
							Name: arg.Name.String(),
							Type: typ,
						}
					}),
					Returns: sliceutil.Mapper(method.Returns, func(ret *ast.Return) MethodReturn {
						typ := ret.Type.String()
						if isMessageType(ret.Type.String()) {
							typ = fmt.Sprintf("*%s", typ)
						}
						if ret.Stream && ret.Type.String() != "[]byte" {
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

	return nil
}

type RpcService struct {
	Name    string
	Methods Methods
}

func (s RpcService) TopicName() string {
	return fmt.Sprintf("TopicRpc%sService", strcase.ToPascal(s.Name))
}

func (s RpcService) TopicValue() string {
	return fmt.Sprintf("ella.rpc.%s.*", strcase.ToSnake(s.Name))
}

type RpcServices []RpcService

func (s *RpcServices) Parse(prog *ast.Program) error {
	isMessageType := createIsMessageTypeFunc(astutil.GetMessages(prog))

	*s = sliceutil.Mapper(astutil.GetServices(prog), func(service *ast.Service) RpcService {
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
						typ := arg.Type.String()
						if isMessageType(arg.Type.String()) {
							typ = fmt.Sprintf("*%s", typ)
						}

						return MethodArg{
							Name: arg.Name.String(),
							Type: typ,
						}
					}),
					Returns: sliceutil.Mapper(method.Returns, func(ret *ast.Return) MethodReturn {
						typ := ret.Type.String()
						if isMessageType(ret.Type.String()) {
							typ = fmt.Sprintf("*%s", typ)
						}

						return MethodReturn{
							Name: ret.Name.String(),
							Type: typ,
						}
					}),
				}
			}),
		}
	})

	return nil
}

func parseMethodOptions(method *ast.Method) MethodOptions {
	mapper := make(map[string]any)
	for _, opt := range method.Options {
		var value any
		switch opt.Value.(type) {
		case *ast.ValueString:
			value = opt.Value.(*ast.ValueString).Value
		case *ast.ValueBool:
			value = opt.Value.(*ast.ValueBool).Value
		case *ast.ValueInt:
			value = opt.Value.(*ast.ValueInt).Value
		case *ast.ValueFloat:
			value = opt.Value.(*ast.ValueFloat).Value
		}
		mapper[strcase.ToPascal(opt.Name.Token.Literal)] = value
	}

	return MethodOptions{
		HttpMethod:    "http." + strcase.ToPascal("Method"+castString(mapper["HttpMethod"], "POST")),
		MaxUploadSize: castInt64(mapper["MaxUploadSize"], 1*1024*1024),
		RawControl:    castBool(mapper["RawControl"], false),
	}
}
