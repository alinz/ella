package parser

import (
	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/token"
	"compiler.ella.to/pkg/strcase"
)

func ParseCustomError(p *Parser) (customError *ast.CustomError, err error) {
	if p.Peek().Type != token.CustomError {
		return nil, p.WithError(p.Peek(), "expected 'error' keyword")
	}

	customError = &ast.CustomError{Token: p.Next()}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a custom error")
	}

	nameTok := p.Next()

	if !strcase.IsPascal(nameTok.Literal) {
		return nil, p.WithError(nameTok, "custom error name must be in Pascal Case format")
	}

	customError.Name = &ast.Identifier{Token: nameTok}

	if p.Peek().Type != token.OpenCurly {
		return nil, p.WithError(p.Peek(), "expected '{' after custom error declaration")
	}

	p.Next() // skip '{'

	// parse Code, HttpStatus and Msg (3 times)
	for p.Peek().Type != token.CloseCurly {
		err = parseCustomErrorValues(p, customError)
		if err != nil {
			return nil, err
		}

	}

	p.Next() // skip '}'

	// check if all values are defined
	if customError.Code == 0 {
		return nil, p.WithError(customError.Token, "code is not defined in custom error")
	}

	if customError.HttpStatus == 0 {
		return nil, p.WithError(customError.Token, "http status is not defined in custom error")
	}

	if customError.Msg == nil {
		return nil, p.WithError(customError.Token, "message is not defined in custom error")
	}

	return customError, nil
}

func parseCustomErrorValues(p *Parser, customError *ast.CustomError) (err error) {
	if p.Peek().Type != token.Identifier {
		return p.WithError(p.Peek(), "expected identifier for defining a custom error value")
	}

	switch p.Peek().Literal {
	case "Code":
		return parseCustomErrorCode(p, customError)
	case "HttpStatus":
		return parseCustomErrorHttpStatus(p, customError)
	case "Msg":
		return parseCustomErrorMsg(p, customError)
	}

	return p.WithError(p.Peek(), "unexpected field name in custom error")
}

func parseCustomErrorCode(p *Parser, customError *ast.CustomError) (err error) {
	if customError.Code != 0 {
		return p.WithError(p.Peek(), "code is already defined in custom error")
	}

	p.Next() // skip 'Code'

	if p.Peek().Type != token.Assign {
		return p.WithError(p.Peek(), "expected '=' after 'Code'")
	}

	p.Next() // skip '='

	if p.Peek().Type != token.ConstInt {
		return p.WithError(p.Peek(), "expected integer value for 'Code'")
	}

	codeValue, err := ParseValue(p)
	if err != nil {
		return err
	}

	customError.Code = codeValue.(*ast.ValueInt).Value

	return nil
}

func parseCustomErrorHttpStatus(p *Parser, customError *ast.CustomError) (err error) {
	if customError.HttpStatus != 0 {
		return p.WithError(p.Peek(), "HttpStatus is already defined in custom error")
	}

	p.Next() // skip 'HttpStatus'

	if p.Peek().Type != token.Assign {
		return p.WithError(p.Peek(), "expected '=' after 'HttpStatus'")
	}

	p.Next() // skip '='

	if p.Peek().Type != token.Identifier {
		return p.WithError(p.Peek(), "expected an HttpStatus value 'HttpStatus'")
	}

	httpStatus, ok := ast.HttpStatusString2Code[p.Peek().Literal]
	if !ok {
		return p.WithError(p.Peek(), "unexpected http status value")
	}

	customError.HttpStatus = httpStatus

	p.Next() // skip http status

	return nil
}

func parseCustomErrorMsg(p *Parser, customError *ast.CustomError) (err error) {
	if customError.Msg != nil {
		return p.WithError(p.Peek(), "Msg is already defined in custom error")
	}

	p.Next() // skip 'Msg'

	if p.Peek().Type != token.Assign {
		return p.WithError(p.Peek(), "expected '=' after 'Msg'")
	}

	p.Next() // skip '='

	msgValue, err := ParseValue(p)
	if err != nil {
		return err
	}

	stringMsgValue, ok := msgValue.(*ast.ValueString)
	if !ok {
		return p.WithError(p.Peek(), "expected string value for 'Msg'")
	}

	customError.Msg = stringMsgValue

	return nil
}
