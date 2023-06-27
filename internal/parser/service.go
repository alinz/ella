package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseService(p *Parser) (service *ast.Service, err error) {
	if p.Peek().Type != token.Service {
		return nil, p.WithError(p.Peek(), "expected service keyword")
	}

	service = &ast.Service{Token: p.Next()}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a service")
	}

	service.Name = &ast.Identifier{Token: p.Next()}

	if p.Peek().Type != token.OpenCurly {
		return nil, p.WithError(p.Peek(), "expected '{' after service declaration")
	}

	p.Next() // skip '{'

	for p.Peek().Type != token.CloseCurly {
		method, err := ParseServiceMethod(p)
		if err != nil {
			return nil, err
		}

		service.Methods = append(service.Methods, method)
	}

	p.Next() // skip '}'

	return service, nil
}

func ParseServiceMethod(p *Parser) (method *ast.Method, err error) {
	switch p.Peek().Type {
	case token.Http, token.Rpc:
	default:
		return nil, p.WithError(p.Peek(), "expected 'http' or 'rpc' keyword")
	}

	method = &ast.Method{
		Type:    p.Next(),
		Args:    make([]*ast.Arg, 0),
		Returns: make([]*ast.Return, 0),
		Options: make([]*ast.Const, 0),
	}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a service method")
	}

	method.Name = &ast.Identifier{Token: p.Next()}

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
	if p.Peek().Type != token.OpenCurly {
		return method, nil
	}

	p.Next() // skip '{'

	for p.Peek().Type != token.CloseCurly {
		option, err := ParseServiceMethodOption(p)
		if err != nil {
			return nil, err
		}

		method.Options = append(method.Options, option)
	}

	p.Next() // skip '}'

	return method, nil
}

func ParseServiceMethodArgument(p *Parser) (arg *ast.Arg, err error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a service method argument")
	}

	arg = &ast.Arg{Name: &ast.Identifier{Token: p.Next()}}

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

	ret = &ast.Return{Name: &ast.Identifier{Token: p.Next()}}

	if p.Peek().Type != token.Colon {
		return nil, p.WithError(p.Peek(), "expected ':' after service method argument name")
	}

	p.Next() // skip ':'

	ret.Type, err = ParseType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type == token.Comma {
		p.Next() // skip ','
	}

	return ret, nil
}

func ParseServiceMethodOption(p *Parser) (constant *ast.Const, err error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a message field constant")
	}

	constant = &ast.Const{
		Name: &ast.Identifier{Token: p.Next()},
	}

	if p.Peek().Type != token.Colon {
		constant.Value = &ast.ValueBool{
			Token:   nil,
			Value:   true,
			Defined: false,
		}

		return constant, nil
	}

	p.Next() // skip ':'

	constant.Value, err = ParseValue(p)
	if err != nil {
		return nil, err
	}

	return constant, nil
}
