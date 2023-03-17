package typescript

import (
	"fmt"

	"ella.to/schema/ast"
	"ella.to/transform"
)

func Signature() transform.Func {
	return func(out transform.Writer) error {
		out.
			Str(`//`).Lines(1).
			Str(`// Code generated by https://ella.to, DO NOT EDIT.`).Lines(1).
			Str(`//`).Lines(1)
		return nil
	}
}

func containsStream(method *ast.Method) bool {
	for _, ret := range method.Returns {
		if ret.Stream {
			return true
		}
	}
	return false
}

func parseType(typ ast.Type) (string, error) {
	switch t := typ.(type) {
	case *ast.TypeBool:
		return `boolean`, nil
	case *ast.TypeInt:
		return `number`, nil
	case *ast.TypeFloat:
		return `number`, nil
	case *ast.TypeUint:
		return `number`, nil
	case *ast.TypeString:
		return `string`, nil
	case *ast.TypeAny:
		return `any`, nil
	case *ast.TypeTimestamp:
		return `string`, nil
	case *ast.TypeArray:
		typ, err := parseType(t.Type)
		if err != nil {
			return "", err
		}
		return typ + "[]", nil
	case *ast.TypeMap:
		key, err := parseType(t.Key)
		if err != nil {
			return "", err
		}
		value, err := parseType(t.Value)
		if err != nil {
			return "", err
		}
		return `{ [key: ` + key + `]: ` + value + ` }`, nil
	case *ast.TypeCustom:
		return t.Name, nil
	default:
		return "", fmt.Errorf("unknown type: %T", t)
	}
}

func getConstByKey(options []*ast.Constant, key string) *ast.Constant {
	for _, option := range options {
		if option.Name.Name == key {
			return option
		}
	}
	return nil
}

func getConstValueAsString(constant *ast.Constant) string {
	if constant == nil {
		return ""
	}

	value, ok := constant.Value.(*ast.ValueString)
	if !ok {
		return ""
	}

	return value.Content
}