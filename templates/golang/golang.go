package golang

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/alinz/rpc.go/pkg/stringcase"
	"github.com/alinz/rpc.go/schema/ast"
)

//go:embed tmpl/*.tmpl
var files embed.FS

func toStructArgs(args []Arg) string {
	var sb strings.Builder

	for i, arg := range args {
		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(stringcase.ToPascal(arg.Name))
		sb.WriteString(" ")
		sb.WriteString(arg.Type)
		sb.WriteString(" ")
		sb.WriteString("`json:\"")
		sb.WriteString(stringcase.ToSnake(arg.Name))
		sb.WriteString("\"`")
	}

	return sb.String()
}

func toStructReturns(returns []Return) string {

	var sb strings.Builder

	for i, ret := range returns {
		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(stringcase.ToPascal(ret.Name))
		sb.WriteString(" ")
		sb.WriteString(ret.Type)
		sb.WriteString(" ")
		sb.WriteString("`json:\"")
		sb.WriteString(stringcase.ToSnake(ret.Name))
		sb.WriteString("\"`")
	}

	return sb.String()

}

var tmplFuncs = template.FuncMap{
	"ToUpper":  strings.ToUpper,
	"ToLower":  strings.ToLower,
	"ToCamel":  stringcase.ToCamel,
	"ToPascal": stringcase.ToPascal,
	"GetArgsLength": func(args []Arg) int {
		return len(args)
	},
	"GetReturnsLength": func(returns []Return) int {
		return len(returns)
	},

	"MethodArgs": func(args []Arg) string {
		var sb strings.Builder

		sb.WriteString("ctx context.Context")

		for _, arg := range args {
			sb.WriteString(", ")
			sb.WriteString(arg.Name)
			sb.WriteString(" ")
			sb.WriteString(arg.Type)
		}

		return sb.String()
	},

	"MethodReturns": func(returns []Return) string {
		var sb strings.Builder

		for i, ret := range returns {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(ret.Name)
			sb.WriteString(" ")
			if ret.Stream {
				sb.WriteString("<-chan ")
			}
			sb.WriteString(ret.Type)
		}

		if sb.Len() > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString("err error")

		return sb.String()
	},

	"IsStream": func(returns []Return) bool {
		for _, ret := range returns {
			if ret.Stream {
				return true
			}
		}

		return false
	},

	"ToStructArgs": toStructArgs,

	"ToStructReturns": toStructReturns,

	"ToExtractArgs": func(args []Arg) string {
		var sb strings.Builder

		for i, arg := range args {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString("args.")
			sb.WriteString(stringcase.ToPascal(arg.Name))
			sb.WriteString(",")
		}

		return sb.String()
	},
	"ToStreamReturnsName": func(returns []Return) string {
		var sb strings.Builder

		for i, ret := range returns {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(ret.Name)
			sb.WriteString("Stream")
		}

		return sb.String()
	},

	"ReturnsExtract": func(returns []Return) string {
		var sb strings.Builder

		for i, ret := range returns {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString("ret.")
			sb.WriteString(stringcase.ToPascal(ret.Name))
		}

		if sb.Len() > 0 {
			sb.WriteString(", err")
		} else {
			sb.WriteString("err")
		}
		return sb.String()
	},

	"MethodArgsStructClient": func(args []Arg) string {
		var sb strings.Builder

		if len(args) > 0 {
			sb.WriteString("struct {\n")
			sb.WriteString(toStructArgs(args))
			sb.WriteString("}{\n")
			for _, arg := range args {
				sb.WriteString(stringcase.ToPascal(arg.Name))
				sb.WriteString(": ")
				sb.WriteString(arg.Name)
				sb.WriteString(",\n")
			}
			sb.WriteString("}\n")
		} else {
			sb.WriteString("emptyStruct{}")
		}

		return sb.String()
	},

	"MethodReturnsStructClient": func(returns []Return) string {
		var sb strings.Builder

		if len(returns) > 0 {
			sb.WriteString("struct {\n")
			sb.WriteString(toStructReturns(returns))
			sb.WriteString("}{}")
		} else {
			sb.WriteString("emptyStruct{}")
		}

		return sb.String()
	},

	"ReturnsOut": func(returns []Return) string {
		var sb strings.Builder

		for i, ret := range returns {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString("out.")
			sb.WriteString(stringcase.ToPascal(ret.Name))
		}

		if sb.Len() > 0 {
			sb.WriteString(", err")
		} else {
			sb.WriteString("err")
		}

		return sb.String()
	},
}

type Constant struct {
	Key   string
	Value any
}

type Constants []Constant

type Enum struct {
	Name      string
	Type      string
	Constants Constants
}

type Enums []Enum

type Field struct {
	Name string
	Type string
	Tags string
}

type Message struct {
	Name   string
	Fields []Field
}

type Arg struct {
	Name string
	Type string
}

type Return struct {
	Name   string
	Type   string
	Stream bool
}

type Method struct {
	Name    string
	Args    []Arg
	Returns []Return
}

type Service struct {
	Name    string
	Methods []Method
}

func generateHeader(out io.Writer, pkgName string) error {
	tmpl, err := loadTemplate("header")
	if err != nil {
		return err
	}

	return tmpl.Execute(out, map[string]string{
		"PackageName": pkgName,
	})
}

func genarateConstants(out io.Writer, nodes []*ast.Constant) error {
	tmpl, err := loadTemplate("constants")
	if err != nil {
		return err
	}

	return tmpl.Execute(out, parseConstants(nodes))
}

func genarateEnums(out io.Writer, nodes []*ast.Enum) error {
	tmpl, err := loadTemplate("enums")
	if err != nil {
		return err
	}

	enums := make(Enums, 0, len(nodes))

	for _, node := range nodes {
		enum := Enum{
			Name:      node.Name.Token.Val,
			Type:      node.Type.TokenLiteral(),
			Constants: parseConstants(node.Constants),
		}

		for i, constant := range enum.Constants {
			enum.Constants[i].Key = fmt.Sprintf("%s_%s", enum.Name, constant.Key)
		}

		enums = append(enums, enum)
	}

	return tmpl.Execute(out, enums)
}

func genarateMessages(out io.Writer, nodes []*ast.Message, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	tmpl, err := loadTemplate("messages")
	if err != nil {
		return err
	}

	var messages []*Message

	// If there are any extends, we add them to the fields first
	for _, node := range nodes {
		msgFields := make([]*ast.Field, 0, len(node.Fields)+len(node.Extends))
		for _, extend := range node.Extends {
			msg, ok := messagesMap[extend.Name]
			if !ok {
				return fmt.Errorf("message %s extended inside %s, not found", extend.Name, node.Name.Token.Val)
			}
			msgFields = append(msgFields, msg.Fields...)
		}
		msgFields = append(msgFields, node.Fields...)

		// now we have all the fields
		message := &Message{
			Name:   node.Name.Token.Val,
			Fields: make([]Field, len(msgFields)),
		}

		for i, field := range msgFields {
			message.Fields[i] = Field{
				Name: field.Name.Token.Val,
				Type: parseFieldType(field.Type, messagesMap, enumsMap),
				Tags: parseFieldOptions(field.Options),
			}
		}

		messages = append(messages, message)
	}

	return tmpl.Execute(out, messages)
}

func genarateServices(out io.Writer, nodes []*ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	tmpl, err := loadTemplate("services")
	if err != nil {
		return err
	}

	services := parseServices(nodes, messagesMap, enumsMap)

	return tmpl.Execute(out, services)
}

func genarateServers(out io.Writer, nodes []*ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	tmpl, err := loadTemplate("servers")
	if err != nil {
		return err
	}

	services := parseServices(nodes, messagesMap, enumsMap)

	return tmpl.Execute(out, services)
}

func genarateClients(out io.Writer, nodes []*ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	tmpl, err := loadTemplate("clients")
	if err != nil {
		return err
	}

	services := parseServices(nodes, messagesMap, enumsMap)

	return tmpl.Execute(out, services)
}

func generateHelper(out io.Writer) error {
	tmpl, err := template.ParseFS(files, "tmpl/helper.gen.tmpl")
	if err != nil {
		return err
	}

	return tmpl.Execute(out, nil)
}

func Generate(out io.Writer, pkgName string, program *ast.Program) error {
	var (
		constants []*ast.Constant
		enums     []*ast.Enum
		messages  []*ast.Message
		services  []*ast.Service
	)

	var (
		enumsMap    = make(map[string]*ast.Enum)
		messagesMap = make(map[string]*ast.Message)
	)

	for _, node := range program.Nodes {
		switch node := node.(type) {
		case *ast.Constant:
			constants = append(constants, node)
		case *ast.Enum:
			enums = append(enums, node)
			enumsMap[node.Name.Token.Val] = node
		case *ast.Message:
			messages = append(messages, node)
			messagesMap[node.Name.Token.Val] = node
		case *ast.Service:
			services = append(services, node)
		}
	}

	err := checkErrs(
		generateHeader(out, pkgName),
		genarateConstants(out, constants),
		genarateEnums(out, enums),
		genarateMessages(out, messages, messagesMap, enumsMap),
		genarateServices(out, services, messagesMap, enumsMap),
		genarateServers(out, services, messagesMap, enumsMap),
		genarateClients(out, services, messagesMap, enumsMap),
		generateHelper(out),
	)

	if err != nil {
		return err
	}

	return nil
}

func checkErrs(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func loadTemplate(name string) (*template.Template, error) {
	file := fmt.Sprintf("%s.gen.tmpl", name)
	tmpl := template.New(file).Funcs(tmplFuncs)
	return tmpl.ParseFS(files, fmt.Sprintf("tmpl/%s", file))
}

func parseConstants(nodes []*ast.Constant) Constants {
	constants := make(Constants, 0, len(nodes))

	for _, node := range nodes {
		constant := Constant{
			Key: stringcase.ToPascal(node.Name.Token.Val),
		}

		switch v := node.Value.(type) {
		case *ast.ValueString:
			constant.Value = fmt.Sprintf(`"%s"`, v.Content)
		case *ast.ValueInt:
			constant.Value = v.Content
		case *ast.ValueUint:
			constant.Value = v.Content
		case *ast.ValueFloat:
			constant.Value = v.Content
		case *ast.ValueBool:
			constant.Value = v.Content
		}

		constants = append(constants, constant)
	}

	return constants
}

func parseFieldOptions(opts []*ast.Constant) string {
	var sb strings.Builder

	mapper := make(map[string]string)
	for _, opt := range opts {
		mapper[opt.Name.Token.Val] = opt.Value.TokenLiteral()
	}

	json, ok := mapper["json"]
	if ok {
		sb.WriteString(`json:"`)
		sb.WriteString(json)
		sb.WriteString(`"`)
	}

	yaml, ok := mapper["yaml"]
	if ok {
		if sb.Len() > 0 {
			sb.WriteString(` `)
		}
		sb.WriteString(`yaml:"`)
		sb.WriteString(yaml)
		sb.WriteString(`"`)
	}

	return sb.String()
}

func parseFieldType(typ ast.Type, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) string {
	switch typ := typ.(type) {
	case *ast.TypeCustom:
		// all the Message type should be pointer
		// all the Enum type should be value
		if _, ok := messagesMap[typ.Name]; ok {
			return fmt.Sprintf("*%s", typ.Name)
		}
		if _, ok := enumsMap[typ.Name]; ok {
			return typ.Name
		}

		panic(fmt.Sprintf("unknown type: %s", typ.Name))
	case *ast.TypeAny:
		return "any"
	case *ast.TypeInt:
		return fmt.Sprintf("int%d", typ.Size)
	case *ast.TypeUint:
		return fmt.Sprintf("uint%d", typ.Size)
	case *ast.TypeByte:
		return "byte"
	case *ast.TypeFloat:
		return fmt.Sprintf("float%d", typ.Size)
	case *ast.TypeString:
		return "string"
	case *ast.TypeBool:
		return "bool"
	case *ast.TypeTimestamp:
		return "time.Time"
	case *ast.TypeMap:
		return fmt.Sprintf("map[%s]%s", parseFieldType(typ.Key, messagesMap, enumsMap), parseFieldType(typ.Value, messagesMap, enumsMap))
	case *ast.TypeArray:
		return fmt.Sprintf("[]%s", parseFieldType(typ.Type, messagesMap, enumsMap))
	}

	panic(fmt.Sprintf("unknown type: %T", typ))
}

func parseServices(nodes []*ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) []*Service {
	services := make([]*Service, 0, len(nodes))

	for _, node := range nodes {
		service := Service{
			Name: stringcase.ToPascal(node.Name.Token.Val),
		}

		for _, method := range node.Methods {
			service.Methods = append(service.Methods, parseMethod(method, messagesMap, enumsMap))
		}

		services = append(services, &service)
	}

	return services
}

func parseMethod(node *ast.Method, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) Method {
	method := Method{
		Name: stringcase.ToPascal(node.Name.Token.Val),
	}

	if node.Args != nil {
		method.Args = parseArgs(node.Args, messagesMap, enumsMap)
	}

	if node.Returns != nil {
		method.Returns = parseReturns(node.Returns, messagesMap, enumsMap)
	}

	return method
}

func parseArgs(nodes []*ast.Arg, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) []Arg {
	var args []Arg

	for _, node := range nodes {
		arg := Arg{
			Name: stringcase.ToCamel(node.Name.Token.Val),
			Type: parseFieldType(node.Type, messagesMap, enumsMap),
		}

		args = append(args, arg)
	}

	return args
}

func parseReturns(nodes []*ast.Return, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) []Return {
	var returns []Return

	for _, node := range nodes {
		ret := Return{
			Name:   stringcase.ToCamel(node.Name.Token.Val),
			Type:   parseFieldType(node.Type, messagesMap, enumsMap),
			Stream: node.Stream,
		}

		returns = append(returns, ret)
	}

	return returns
}
