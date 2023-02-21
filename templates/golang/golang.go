package golang

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/alinz/rpc.go/schema/ast"
)

//go:embed tmpl/*.tmpl
var files embed.FS

var tmplFuncs = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
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
	return nil
}

func genarateServers(out io.Writer, nodes []*ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	return nil
}

func genarateClients(out io.Writer, nodes []*ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	return nil
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
			Key: node.Name.Token.Val,
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

	return ""
}
