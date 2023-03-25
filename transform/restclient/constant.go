package restclient

import (
	"ella.to/schema/ast"
	"ella.to/transform"
)

func Constants(program *ast.Program) transform.Func {
	return func(out transform.Writer) error {
		constants := ast.GetSlice[*ast.Constant](program)

		isHostnameDefined := false
		isPortDefined := false

		for _, constant := range constants {
			switch constant.Name.Name {
			case "http_hostname":
				out.Str("@hostname = %s", constant.Value.(*ast.ValueString).Content).Lines(1)
				isHostnameDefined = true
			case "http_port":
				out.Str("@port = %s", constant.Value.(*ast.ValueString).Content).Lines(1)
				isPortDefined = true
			}
		}

		if !isHostnameDefined {
			out.Str("@hostname = ").Lines(1)
		}
		if !isPortDefined {
			out.Str("@port = ").Lines(1)
		}

		out.Str("@host = {{hostname}}:{{port}}").Lines(1)
		out.Str("@authorization = ")

		return nil
	}
}
