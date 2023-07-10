package golang

import (
	"io"
	"sort"

	"ella.to/internal/ast"
	"ella.to/internal/code"
)

func generateEnums(out io.Writer, enums []*ast.Enum) error {
	tmpl, err := code.LoadTemplate(files, "enums")
	if err != nil {
		return err
	}

	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name.String() < enums[j].Name.String()
	})

	templEnums := code.Mapper(enums, func(enum *ast.Enum) any {
		keys := code.Mapper(enum.Constants, func(key *ast.Const) any {
			return struct {
				Name  string
				Value string
			}{
				Name:  key.Name.String(),
				Value: key.Value.String(),
			}
		})

		return struct {
			Name string
			Type string
			Keys []any
		}{
			Name: enum.Name.String(),
			Type: enum.Type.String(),
			Keys: keys,
		}
	})

	return tmpl.Execute(out, templEnums)
}
