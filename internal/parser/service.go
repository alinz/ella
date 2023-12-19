package parser

import (
	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/token"
	"compiler.ella.to/pkg/strcase"
)

func ParseService(p *Parser) (service *ast.Service, err error) {
	if p.Peek().Type != token.Service {
		return nil, p.WithError(p.Peek(), "expected service keyword")
	}

	service = &ast.Service{Token: p.Next()}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a service")
	}

	nameTok := p.Next()

	if !strcase.IsPascal(nameTok.Literal) {
		return nil, p.WithError(nameTok, "service name must be in PascalCase format")
	}

	service.Name = &ast.Identifier{Token: nameTok}

	if p.Peek().Type != token.OpenCurly {
		return nil, p.WithError(p.Peek(), "expected '{' after service declaration")
	}

	p.Next() // skip '{'

	for p.Peek().Type != token.CloseCurly {
		methods, err := ParseServiceMethod(p)
		if err != nil {
			return nil, err
		}

		service.Methods = append(service.Methods, methods...)
	}

	p.Next() // skip '}'

	return service, nil
}

func ParseMethodTypes(p *Parser) (methodTypes []ast.MethodType, err error) {
	methodTypes = make([]ast.MethodType, 0)
	done := false
	for !done {
		switch p.Peek().Type {
		case token.Http:
			methodTypes = append(methodTypes, ast.MethodHTTP)
			p.Next() // skip http
		case token.Rpc:
			methodTypes = append(methodTypes, ast.MethodRPC)
			p.Next() // skip rpc
		case token.Comma:
			if len(methodTypes) == 0 {
				return nil, p.WithError(p.Peek(), "expected 'http' or 'rpc' keyword")
			} else if len(methodTypes) == 2 {
				return nil, p.WithError(p.Peek(), "there should be only two method types")
			}

			p.Next() // skip ','
			continue
		default:
			if len(methodTypes) == 0 {
				return nil, p.WithError(p.Peek(), "expected 'http' or 'rpc' keyword")
			} else if len(methodTypes) > 0 {
				done = true
			}
		}
	}

	return methodTypes, nil
}

func ParseServiceMethod(p *Parser) (methods []*ast.Method, err error) {
	methodTypes, err := ParseMethodTypes(p)
	if err != nil {
		return nil, err
	}

	methods = make([]*ast.Method, 0, len(methodTypes))

	method := &ast.Method{
		Args:    make([]*ast.Arg, 0),
		Returns: make([]*ast.Return, 0),
		Options: make([]*ast.Option, 0),
	}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a service method")
	}

	nameTok := p.Next()

	if !strcase.IsPascal(nameTok.Literal) {
		return nil, p.WithError(nameTok, "service method name must be in PascalCase format")
	}

	method.Name = &ast.Identifier{Token: nameTok}

	if p.Peek().Type != token.OpenParen {
		return nil, p.WithError(p.Peek(), "expected '(' after service method name")
	}

	p.Next() // skip '('

	for p.Peek().Type != token.CloseParen {
		arg, err := ParseServiceMethodArgument(p)
		if err != nil {
			return nil, err
		}

		method.Args = append(method.Args, arg)
	}

	p.Next() // skip ')'

	if p.Peek().Type == token.Return {
		p.Next() // skip =>

		if p.Peek().Type != token.OpenParen {
			return nil, p.WithError(p.Peek(), "expected '(' after '=>'")
		}

		p.Next() // skip '('

		for p.Peek().Type != token.CloseParen {
			ret, err := ParseServiceMethodReturnArg(p)
			if err != nil {
				return nil, err
			}

			method.Returns = append(method.Returns, ret)
		}

		p.Next() // skip ')'
	}

	// we return early if there are no options
	// as options are defined by curly braces
	if p.Peek().Type == token.OpenCurly {
		method.Options, err = ParseOptions(p)
		if err != nil {
			return nil, err
		}
	}

	for _, methodType := range methodTypes {
		methods = append(methods, &ast.Method{
			Type:    methodType,
			Name:    method.Name,
			Args:    method.Args,
			Returns: method.Returns,
			Options: method.Options,
		})
	}

	return methods, nil
}

func ParseServiceMethodArgument(p *Parser) (arg *ast.Arg, err error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a service method argument")
	}

	nameTok := p.Next()

	if !strcase.IsCamel(nameTok.Literal) {
		return nil, p.WithError(nameTok, "service method argument name must be in camelCase format")
	}

	arg = &ast.Arg{Name: &ast.Identifier{Token: nameTok}}

	if p.Peek().Type != token.Colon {
		return nil, p.WithError(p.Peek(), "expected ':' after service method argument name")
	}

	p.Next() // skip ':'

	arg.Type, err = ParseType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type == token.Comma {
		p.Next() // skip ','
	}

	return arg, nil
}

func ParseServiceMethodReturnArg(p *Parser) (ret *ast.Return, err error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a service method argument")
	}

	nameTok := p.Next()

	if !strcase.IsCamel(nameTok.Literal) {
		return nil, p.WithError(nameTok, "service method argument name must be in camelCase format")
	}

	ret = &ast.Return{Name: &ast.Identifier{Token: nameTok}}

	if p.Peek().Type != token.Colon {
		return nil, p.WithError(p.Peek(), "expected ':' after service method argument name")
	}

	p.Next() // skip ':'

	if p.Peek().Type == token.Stream {
		ret.Stream = true
		p.Next() // skip 'stream'
	}

	ret.Type, err = ParseType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type == token.Comma {
		if ret.Stream {
			return nil, p.WithError(p.Peek(), "there should be only one stream on the return type")
		}
		p.Next() // skip ','
	}

	return ret, nil
}
