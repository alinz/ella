package transform

import (
	"fmt"
	"io"
	"strings"

	"ella.to/pkg/stringcase"
)

type Writer interface {
	Lines(n int) Writer
	Tabs(n int) Writer
	Snake(s string) Writer
	Pascal(s string) Writer
	Camel(s string) Writer
	Str(format string, args ...any) Writer
	StrCond(cond bool, format string, args ...any) Writer
}

type Func func(out Writer) error

func Run(out io.Writer, funcs ...Func) error {
	for _, f := range funcs {
		w := &writer{out: out}

		if err := f(w); err != nil {
			return err
		}

		if err := w.Done(); err != nil {
			return err
		}
	}

	return nil
}

type writer struct {
	err error
	out io.Writer
}

var _ Writer = (*writer)(nil)

func (w *writer) Lines(n int) Writer {
	if w.err != nil {
		return w
	}

	fmt.Fprintf(w.out, "%s", strings.Repeat("\n", n))
	return w
}

func (w *writer) Tabs(n int) Writer {
	if w.err != nil {
		return w
	}

	fmt.Fprintf(w.out, "%s", strings.Repeat("\t", n))
	return w
}

func (w *writer) Snake(s string) Writer {
	return w.Str(stringcase.ToSnake(s))
}

func (w *writer) Pascal(s string) Writer {
	return w.Str(stringcase.ToPascal(s))
}

func (w *writer) Camel(s string) Writer {
	return w.Str(stringcase.ToCamel(s))
}

func (w *writer) Str(format string, args ...any) Writer {
	if w.err != nil {
		return w
	}

	fmt.Fprintf(w.out, format, args...)
	return w
}

func (w *writer) StrCond(cond bool, format string, args ...any) Writer {
	if cond {
		return w.Str(format, args...)
	}

	return w
}

func (w *writer) Done() error {
	return w.err
}
