package http

import (
	"fmt"

	"ella.to/pkg/stringcase"
	"ella.to/schema/ast"
	"ella.to/transform"
)

func getEnumsMap(program *ast.Program) map[string]*ast.Enum {
	enums := ast.GetSlice[*ast.Enum](program)
	enumMap := make(map[string]*ast.Enum)

	for _, enum := range enums {
		enumMap[enum.Name.Name] = enum
	}

	return enumMap
}

func getMessagesMap(program *ast.Program) map[string]*ast.Message {
	messages := ast.GetSlice[*ast.Message](program)
	messagesMap := make(map[string]*ast.Message)

	for _, message := range messages {
		messagesMap[message.Name.Name] = message
	}

	return messagesMap
}

func Services(program *ast.Program) transform.Func {
	return func(out transform.Writer) error {
		services := ast.GetSlice[*ast.Service](program)
		enumsMap := getEnumsMap(program)
		customTypesMap := getMessagesMap(program)

		for _, service := range services {
			serviceName := service.Name.Name

			if len(service.Methods) == 0 {
				continue
			}

			out.Lines(2).Str("### ").Str(serviceName).Lines(1)

			for i, method := range service.Methods {

				if i == 0 {
					out.Lines(1)
				} else {
					out.Lines(2)
				}

				out.Str("### ").Str(serviceName).Str(" :: ").Str(method.Name.Name).Lines(2)

				out.
					Str("POST {{host}}/").Pascal(serviceName).Str("/").Pascal(method.Name.Name).Str(" HTTP/1.1").
					Lines(1)

				out.Str("Content-Type: application/json").Lines(1)
				out.Str("Authorization: BEARER {{authorization}}").Lines(2)

				// get keys as array of strings
				// get the values as array of ast.Type
				// tab is the number of indentations

				keys := make([]string, 0)
				values := make([]ast.Type, 0)
				for _, arg := range method.Args {
					keys = append(keys, arg.Name.Name)
					values = append(values, arg.Type)
				}

				Payload(out, keys, values, 0, enumsMap, customTypesMap)
			}
		}

		return nil
	}
}

func Payload(out transform.Writer, keys []string, values []ast.Type, tab int, enumsMap map[string]*ast.Enum, customTypesMap map[string]*ast.Message) error {
	out.Str(`{`).Lines(1)

	for i, key := range keys {
		lastItem := len(keys) == i+1

		out.Tabs(tab + 1).Str(`"`).Str(stringcase.ToSnake(key)).Str(`": `)
		if err := Type(out, values[i], tab, enumsMap, customTypesMap); err != nil {
			return err
		}

		if !lastItem {
			out.Str(`,`)
		}

		// add comment for extra info
		ExtraTypoInfo(out, values[i], enumsMap)

		out.Lines(1)
	}

	out.Tabs(tab).Str(`}`)
	return nil
}

func ExtraTypoInfo(out transform.Writer, typ ast.Type, enumsMap map[string]*ast.Enum) {
	switch t := typ.(type) {
	case *ast.TypeCustom:
		if enum, ok := enumsMap[t.Name]; ok {
			out.Str(" // ")
			for i, constant := range enum.Constants {
				if i != 0 {
					out.Str(", ")
				}
				out.Str(constant.Name.Name)
			}
		}
	}
}

func Type(out transform.Writer, typ ast.Type, tab int, enumsMap map[string]*ast.Enum, messagesMap map[string]*ast.Message) error {
	switch t := typ.(type) {
	case *ast.TypeBool:
		out.Str("false")
	case *ast.TypeInt:
		out.Str("-0")
	case *ast.TypeFloat:
		out.Str("0.0")
	case *ast.TypeUint:
		out.Str("0")
	case *ast.TypeString:
		out.Str(`""`)
	case *ast.TypeAny:
		out.Str("null")
	case *ast.TypeTimestamp:
		out.Str(`"2021-05-31T00:00:00Z"`)
	case *ast.TypeArray:
		out.Str(`[]`)
	case *ast.TypeMap:
		out.Str(`{}`)
	case *ast.TypeCustom:
		if _, ok := enumsMap[t.Name]; ok {
			out.Str(`""`)
		} else if msg, ok := messagesMap[t.Name]; ok {
			keys := make([]string, 0)
			values := make([]ast.Type, 0)
			for _, field := range getAllMessageFileds(msg, messagesMap) {
				keys = append(keys, field.Name.Name)
				values = append(values, field.Type)
			}

			Payload(out, keys, values, tab+1, enumsMap, messagesMap)
		}

	default:
		return fmt.Errorf("unknown type: %T", typ)
	}

	return nil
}

func getAllMessageFileds(msg *ast.Message, messagesMap map[string]*ast.Message) []*ast.Field {
	fields := make([]*ast.Field, 0)

	for _, extend := range msg.Extends {
		if msg, ok := messagesMap[extend.Name]; ok {
			fields = getAllMessageFieldsHelper(msg, fields)
		}
	}

	return getAllMessageFieldsHelper(msg, fields)
}

func getAllMessageFieldsHelper(msg *ast.Message, fields []*ast.Field) []*ast.Field {
	for _, field := range msg.Fields {
		fields = append(fields, field)
	}

	return fields
}
