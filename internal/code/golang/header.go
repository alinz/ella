package golang

import (
	"io"

	"ella.to/internal/code"
)

func generateHeader(out io.Writer, pkg string) error {
	tmpl, err := code.LoadTemplate(files, "header")
	if err != nil {
		return err
	}

	return tmpl.Execute(out, struct {
		PkgName string
	}{
		PkgName: pkg,
	})
}
