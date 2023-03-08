package code_test

import (
	"bytes"
	"io"
	"testing"

	"ella.to/pkg/code"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	const src = `
	package plugin

	func Fib(n int) int {
		return fib(n, 0, 1)
	}
	
	func fib(n, a, b int) int {
		if n == 0 {
			return a
		} else if n == 1 {
			return b
		}
		return fib(n-1, b, a+b)
	}	
	`

	fib, err := code.Func[func(int) int](src, "plugin.Fib")
	assert.NoError(t, err)
	assert.NotNil(t, fib)
	assert.Equal(t, fib(10), 55)
}

type Point struct {
	X int
	Y int
}

func (p *Point) Sum() int {
	return p.X + p.Y
}

func TestCustomStruct(t *testing.T) {
	const src = `
	package test

	import "example.com/cool"

	func sum(p *cool.Point) int {
		return p.X + p.Y
	}

	func Sum() int {
		point := &cool.Point{
			X: 1,
			Y: 2,
		}

		return sum(point)
	}
	

	func Example(p *cool.Point) int {
		return p.Sum()
	}
	`

	fn, err := code.Func[func() int](src, "test.Sum", code.WithExposeStuct[Point]("example.com/cool", "Point"))
	assert.NoError(t, err)
	assert.NotNil(t, fn)
	assert.Equal(t, fn(), 3)

	fn2, err := code.Func[func(*Point) int](src, "test.Example",
		code.WithExposeStuct[Point]("example.com/cool", "Point"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, fn2)
	assert.Equal(t, fn2(&Point{X: 1, Y: 2}), 3)
}

func TestSharedWriter(t *testing.T) {
	const src = `
	package plugin

	import "io"

	func Write(w io.Writer, s string) (int, error) {
		return w.Write([]byte(s))
	}
	`

	var buffer bytes.Buffer

	fn, err := code.Func[func(w io.Writer, value string) (int, error)](src, "plugin.Write")
	assert.NoError(t, err)
	assert.NotNil(t, fn)
	n, err := fn(&buffer, "abc")
	assert.NoError(t, err)
	assert.Equal(t, n, 3)
	assert.Equal(t, buffer.String(), "abc")
}
