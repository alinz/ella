package golang

import (
	"io"

	"ella.to/internal/code"
)

func generateHelper(out io.Writer) error {
	tmpl, err := code.LoadTemplate(files, "helper")
	if err != nil {
		return err
	}

	return tmpl.Execute(out, struct {
	}{})
}
