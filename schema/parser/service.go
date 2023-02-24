package parser

import (
	"fmt"

	"github.com/alinz/rpc.go/schema/ast"
	"github.com/alinz/rpc.go/schema/token"
)

func (p *Parser) parseService() (*ast.Service, error) {
	if p.nextToken.Kind != token.Service {
		return nil, fmt.Errorf("expected service keyword but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip service keyword

	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected identifier but got %s", p.nextToken.Kind)
	}

	service := &ast.Service{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
	}

	p.scanToken() // skip service name

	if p.nextToken.Kind != token.OpenCurl {
		return nil, fmt.Errorf("expected { but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip {

	for p.nextToken.Kind != token.CloseCurl {
		method, err := p.parseMethod()
		if err != nil {
			return nil, err
		}

		service.Methods = append(service.Methods, method)
	}

	p.scanToken() // skip }

	return service, nil
}

func (p *Parser) parseMethod() (*ast.Method, error) {
	var err error

	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected identifier but got %s", p.nextToken.Kind)
	}

	method := &ast.Method{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
		Options: make([]*ast.Constant, 0),
		Args:    make([]*ast.Arg, 0),
		Returns: make([]*ast.Return, 0),
	}

	p.scanToken() // skip method name

	if p.nextToken.Kind != token.OpenParen {
		return nil, fmt.Errorf("expected ( but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip (

	for p.nextToken.Kind != token.CloseParen {
		arg, err := p.parseArg()
		if err != nil {
			return nil, err
		}

		method.Args = append(method.Args, arg)
	}

	p.scanToken() // skip )

	if p.nextToken.Kind == token.OpenCurl {
		method.Options, err = p.parseMethodOptions()
		if err != nil {
			return nil, err
		}
	}

	if p.nextToken.Kind != token.Return {
		return method, nil
	}

	p.scanToken() // skip =>

	if p.nextToken.Kind != token.OpenParen {
		return nil, fmt.Errorf("expected ( for return args but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip (

	for p.nextToken.Kind != token.CloseParen {
		ret, err := p.parseReturn()
		if err != nil {
			return nil, err
		}

		method.Returns = append(method.Returns, ret)
	}

	p.scanToken() // skip )

	if p.nextToken.Kind == token.OpenCurl {
		method.Options, err = p.parseMethodOptions()
		if err != nil {
			return nil, err
		}
	}

	return method, nil
}

func (p *Parser) parseMethodOptions() (options []*ast.Constant, err error) {
	if p.nextToken.Kind != token.OpenCurl {
		return nil, fmt.Errorf("expected { but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip {

	for p.nextToken.Kind != token.CloseCurl {
		option, err := p.parseConstant(true)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	p.scanToken() // skip }

	return options, nil
}

func (p *Parser) parseArg() (arg *ast.Arg, err error) {
	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected identifier but got %s", p.nextToken.Kind)
	}

	arg = &ast.Arg{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
	}

	p.scanToken() // skip arg name

	if p.nextToken.Kind != token.Colon {
		return nil, fmt.Errorf("expected : but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip :

	arg.Type, err = p.parseType()
	if err != nil {
		return nil, err
	}

	if p.nextToken.Kind == token.Comma {
		p.scanToken() // skip ,
	}

	return arg, nil
}

func (p *Parser) parseReturn() (ret *ast.Return, err error) {
	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected identifier but got %s", p.nextToken.Kind)
	}

	ret = &ast.Return{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
	}

	p.scanToken() // skip return name

	if p.nextToken.Kind != token.Colon {
		return nil, fmt.Errorf("expected : but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip :

	if p.nextToken.Kind == token.Stream {
		ret.Stream = true
		p.scanToken() // skip stream
	}

	ret.Type, err = p.parseType()
	if err != nil {
		return nil, err
	}

	return ret, nil
}
