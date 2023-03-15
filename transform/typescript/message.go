package typescript

import (
	"ella.to/schema/ast"
	"ella.to/transform"
)

func Message(message *ast.Message) transform.Func {
	return func(out transform.Writer) error {
		out.
			Str(`export interface `).Pascal(message.Name.Name).Str(` {`).
			Lines(1)

		for _, field := range message.Fields {
			var customName string

			json := getConstByKey(field.Options, "json")
			if json != nil {
				customName = getConstValueAsString(json)
			}

			typ, err := parseType(field.Type)
			if err != nil {
				return err
			}

			out.Tabs(1)
			if customName != "" {
				out.Str(customName)
			} else {
				out.Snake(field.Name.Name)
			}
			out.Str(`: `).Str(typ).Lines(1)
		}

		out.Str(`}`)

		return nil
	}
}

func Messages(messages []*ast.Message) transform.Func {
	return func(out transform.Writer) error {
		out.Str("// Messages").Lines(2)

		for _, message := range messages {
			if err := Message(message)(out); err != nil {
				return err
			}
			out.Lines(2)
		}

		return nil
	}
}
