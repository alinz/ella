package http

import (
	"ella.to/schema/ast"
	"ella.to/transform"
)

func Constants(program *ast.Program) transform.Func {
	return func(out transform.Writer) error {
		constants := ast.GetSlice[*ast.Constant](program)

		for _, constant := range constants {
			switch constant.Name.Name {
			case "http_hostname":
				out.Str("@hostname = %s", constant.Value.(*ast.ValueString).Content).Lines(1)
			case "http_port":
				out.Str("@port = %s", constant.Value.(*ast.ValueString).Content).Lines(1)
			}
		}
		out.Str("@host = {{hostname}}:{{port}}").Lines(1)
		out.Str("@authorization = ")

		return nil
	}
}
