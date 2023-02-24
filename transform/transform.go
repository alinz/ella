package transform

import (
	"fmt"
	"io"
)

type Writer interface {
	Tabs(n int) Writer
	NewLines(n int) Writer
	Indents(n int) Writer
	String(format string, args ...any) Writer
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
	err    error
	out    io.Writer
	intend int
}

var _ Writer = (*writer)(nil)

func (w *writer) action(value string) *writer {
	if w.err != nil {
		return w
	}

	_, w.err = w.out.Write([]byte(value))

	return w
}

func (w *writer) actionN(value string, n int) *writer {
	if w.err != nil {
		return w
	}

	for i := 0; i < n; i++ {
		_, w.err = w.out.Write([]byte(value))
		if w.err != nil {
			return w
		}
	}

	return w
}

func (w *writer) Tabs(n int) Writer {
	return w.actionN("\t", n)
}

func (w *writer) NewLines(n int) Writer {
	return w.actionN("\n", n)
}

func (w *writer) Indents(n int) Writer {
	if w.err != nil {
		return w
	}

	w.intend += n

	return w
}

func (w *writer) String(value string, args ...any) Writer {
	w.actionN("\t", w.intend)
	w.action(fmt.Sprintf(value, args...))
	return w
}

func (w *writer) Done() error {
	return w.err
}
