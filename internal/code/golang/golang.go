package golang

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
	"ella.to/internal/token"
	"ella.to/pkg/strcase"
)

//go:embed templates/*.tmpl
var files embed.FS

type Constant struct {
	Name  string
	Value string
}

type Constants []Constant

func (c *Constants) Parse(prog *ast.Program) error {
	*c = code.Mapper(code.GetConstants(prog), func(constant *ast.Const) Constant {
		var value string

		switch v := constant.Value.(type) {
		case *ast.ValueString:
			if v.Token.Type == token.ConstStringSingleQuote {
				value = fmt.Sprintf(`"%s"`, strings.ReplaceAll(v.Token.Val, `"`, `\"`))
			} else {
				value = constant.Value.String()
			}
		case *ast.ValueInt:
			value = strconv.FormatInt(v.Value, 10)
		default:
			value = constant.Value.String()
		}

		return Constant{
			Name:  constant.Name.String(),
			Value: value,
		}
	})

	return nil
}

type EnumKeyValue struct {
	Name  string
	Value string
}

type Enum struct {
	Name string
	Type string
	Keys []EnumKeyValue
}

type Enums []Enum

func (e *Enums) Parse(prog *ast.Program) error {
	*e = code.Mapper(code.GetEnums(prog), func(enum *ast.Enum) Enum {
		return Enum{
			Name: enum.Name.String(),
			Type: enum.Type.String(),
			Keys: code.Mapper(enum.Constants, func(key *ast.Const) EnumKeyValue {
				return EnumKeyValue{
					Name:  key.Name.String(),
					Value: key.Value.String(),
				}
			}),
		}
	})

	return nil
}

type MessageField struct {
	Name string
	Type string
	Tags string
}

type MessageFields []MessageField

func (m *MessageFields) Parse(message *ast.Message, isMessageType func(value string) bool) error {
	*m = code.Mapper(message.Fields, func(field *ast.Field) MessageField {
		typ := parseMessageFiledType(field.Type)
		if isMessageType(field.Type.String()) {
			typ = fmt.Sprintf("*%s", typ)
		}

		return MessageField{
			Name: field.Name.String(),
			Type: typ,
			Tags: parseMessageFieldOptions(field),
		}
	})
	return nil
}

type Message struct {
	Name   string
	Fields MessageFields
}

type Messages []Message

func (m *Messages) Parse(prog *ast.Program) error {
	messages, isMessageType := createIsMessageType(prog)

	*m = code.Mapper(messages, func(message *ast.Message) Message {
		msg := Message{
			Name: message.Name.String(),
		}

		msg.Fields.Parse(message, isMessageType)

		return msg
	})

	return nil
}

type MethodArg struct {
	Name string
	Type string
}

type MethodArgs []MethodArg

func (m MethodArgs) Definitions() string {
	return strings.Join(code.Mapper(m, func(arg MethodArg) string {
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

	return strings.Join(code.Mapper(m, func(arg MethodReturn) string {
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
	return strings.Join(code.Mapper(code.Filter(m.Args, func(arg MethodArg) bool {
		return arg.Type != "<-chan *fileUpload"
	}), func(arg MethodArg) string {
		return prefix + strcase.ToPascal(arg.Name) + ","
	}), "\n")
}

func (m Method) ArgsStructDefinitions() string {
	return strings.Join(code.Mapper(code.Filter(m.Args, func(arg MethodArg) bool {
		return arg.Type != "<-chan *fileUpload"
	}), func(arg MethodArg) string {
		return fmt.Sprintf("%s %s", strcase.ToPascal(arg.Name), arg.Type)
	}), "\n")
}

func (m Method) ArgsNamesValues() string {
	return strings.Join(code.Mapper(m.Args, func(arg MethodArg) string {
		return strcase.ToPascal(arg.Name) + ":" + arg.Name + ","
	}), "\n")
}

func (m Method) ReturnsNames(prefix string) string {
	return strings.Join(code.Mapper(m.Returns, func(arg MethodReturn) string {
		return prefix + strcase.ToPascal(arg.Name) + ", "
	}), "")
}

func (m Method) ReturnsStructDefinitions() string {
	return strings.Join(code.Mapper(m.Returns, func(arg MethodReturn) string {
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
	_, isMessageType := createIsMessageType(prog)

	*s = code.Mapper(code.GetServices(prog), func(service *ast.Service) HttpService {
		methods := code.Filter(service.Methods, func(method *ast.Method) bool {
			return method.Type.Val == "http"
		})

		return HttpService{
			Name: service.Name.String(),
			Methods: code.Mapper(methods, func(method *ast.Method) Method {
				return Method{
					Name:    method.Name.String(),
					Service: service.Name.String(),
					Options: parseMethodOptions(method),
					Args: code.Mapper(method.Args, func(arg *ast.Arg) MethodArg {
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
					Returns: code.Mapper(method.Returns, func(ret *ast.Return) MethodReturn {
						typ := ret.Type.String()
						if isMessageType(ret.Type.String()) {
							typ = fmt.Sprintf("*%s", typ)
						}
						if ret.Stream {
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
	_, isMessageType := createIsMessageType(prog)

	*s = code.Mapper(code.GetServices(prog), func(service *ast.Service) RpcService {
		methods := code.Filter(service.Methods, func(method *ast.Method) bool {
			return method.Type.Val == "rpc"
		})

		return RpcService{
			Name: service.Name.String(),
			Methods: code.Mapper(methods, func(method *ast.Method) Method {
				return Method{
					Name:    method.Name.String(),
					Service: service.Name.String(),
					Args: code.Mapper(method.Args, func(arg *ast.Arg) MethodArg {
						typ := arg.Type.String()
						if isMessageType(arg.Type.String()) {
							typ = fmt.Sprintf("*%s", typ)
						}

						return MethodArg{
							Name: arg.Name.String(),
							Type: typ,
						}
					}),
					Returns: code.Mapper(method.Returns, func(ret *ast.Return) MethodReturn {
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

func parseMessageFiledType(typ ast.Type) string {
	switch typ := typ.(type) {
	case *ast.CustomType:
		return typ.String()
	case *ast.Any:
		return "any"
	case *ast.Int:
		return fmt.Sprintf("int%d", typ.Size)
	case *ast.Uint:
		return fmt.Sprintf("uint%d", typ.Size)
	case *ast.Byte:
		return "byte"
	case *ast.Float:
		return fmt.Sprintf("float%d", typ.Size)
	case *ast.String:
		return "string"
	case *ast.Bool:
		return "bool"
	case *ast.Timestamp:
		return "time.Time"
	case *ast.Map:
		return fmt.Sprintf("map[%s]%s", parseMessageFiledType(typ.Key), parseMessageFiledType(typ.Value))
	case *ast.Array:
		return fmt.Sprintf("[]%s", parseMessageFiledType(typ.Type))
	}

	// This shouldn't happen as the validator should catch this any errors
	panic(fmt.Sprintf("unknown type: %T", typ))
}

func parseMessageFieldOptions(field *ast.Field) string {
	var sb strings.Builder

	mapper := make(map[string]ast.Value)
	for _, opt := range field.Options {
		mapper[strings.ToLower(opt.Name.Token.Val)] = opt.Value
	}

	jsonTagValue := strings.ToLower(strcase.ToSnake(field.Name.String()))

	jsonValue, ok := mapper["json"]
	if ok {
		switch jsonValue := jsonValue.(type) {
		case *ast.ValueString:
			jsonTagValue = jsonValue.Token.Val
		case *ast.ValueBool:
			if !jsonValue.Value {
				jsonTagValue = "-"
			}
		}
	}

	sb.WriteString(`json:"`)
	sb.WriteString(jsonTagValue)
	sb.WriteString(`"`)

	yamlTagValue := strings.ToLower(strcase.ToSnake(field.Name.String()))

	yamlValue, ok := mapper["yaml"]
	if ok {
		switch yamlValue := yamlValue.(type) {
		case *ast.ValueString:
			yamlTagValue = yamlValue.Token.Val
		case *ast.ValueBool:
			if !yamlValue.Value {
				yamlTagValue = "-"
			}
		}
	}

	sb.WriteString(` yaml:"`)
	sb.WriteString(yamlTagValue)
	sb.WriteString(`"`)

	return sb.String()
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
		mapper[strcase.ToPascal(opt.Name.Token.Val)] = value
	}

	return MethodOptions{
		HttpMethod:    "http." + strcase.ToPascal("Method"+castString(mapper["HttpMethod"], "POST")),
		MaxUploadSize: castInt64(mapper["MaxUploadSize"], 1*1024*1024),
		RawControl:    castBool(mapper["RawControl"], false),
	}
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

func createIsMessageType(prog *ast.Program) ([]*ast.Message, func(value string) bool) {
	messages := code.GetMessages(prog)
	messagesMap := make(map[string]struct{})
	for _, message := range messages {
		messagesMap[message.Name.String()] = struct{}{}
	}

	return messages, func(value string) bool {
		_, ok := messagesMap[value]
		return ok
	}
}
