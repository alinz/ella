package code

import (
	"errors"
	"net/url"
	"path"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var (
	ErrCodeSign             = errors.New("code sign is incorrect")
	ErrDuplicatedExportName = errors.New("duplicated export name")
	ErrFuncNotFound         = errors.New("function not found")
)

type Options struct {
	pkgs     []string
	env      []string // example=1
	exported map[string]map[string]reflect.Value
}

type OptionsFunc func(opt *Options) error

func WithPkgs(pkgs ...string) OptionsFunc {
	return func(opt *Options) error {
		opt.pkgs = pkgs
		return nil
	}
}

func WithEnv(env ...string) OptionsFunc {
	return func(opt *Options) error {
		opt.env = env
		return nil
	}
}

func WithExposeStuct[T any](pkg string, name string) OptionsFunc {
	return func(opt *Options) error {
		if opt.exported == nil {
			opt.exported = make(map[string]map[string]reflect.Value)
		}

		// NOTE: in order to make this working we need
		// to grab the pkg name and duplicate the last part of it
		// e.g. github.com/a/b -> github.com/a/b/b
		// we need this done to get the yaegi happy ¯\_(ツ)_/¯
		p, err := url.Parse(pkg)
		if err != nil {
			return err
		}
		p.Path = path.Join(p.Path, path.Base(p.Path))
		pkg = p.String()

		// accessing the pkg name
		ptr, ok := opt.exported[pkg]
		if !ok {
			ptr = make(map[string]reflect.Value)
			opt.exported[pkg] = ptr
		}

		if _, ok = ptr[name]; ok {
			return ErrDuplicatedExportName
		}

		ptr[name] = reflect.ValueOf((*T)(nil))

		return nil
	}
}

func Func[T any](src string, name string, optFns ...OptionsFunc) (T, error) {
	var fn T

	opts := Options{}
	for _, optFn := range optFns {
		if err := optFn(&opts); err != nil {
			return fn, err
		}
	}

	intp := interp.New(interp.Options{
		GoPath: "./_gopath/",
		Env:    opts.env,
	})
	if err := intp.Use(stdlib.Symbols); err != nil {
		return fn, err
	}

	if opts.exported != nil {
		if err := intp.Use(opts.exported); err != nil {
			return fn, err
		}
	}

	if _, err := intp.Eval(src); err != nil {
		return fn, err
	}

	v, err := intp.Eval(name) // "plugin.Fib"
	if err != nil {
		return fn, err
	}

	inter := v.Interface()

	fn, ok := inter.(T)
	if !ok {
		return fn, ErrFuncNotFound
	}

	return fn, nil
}
