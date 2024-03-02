package golang

import (
	"fmt"

	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/ast/astutil"
	"compiler.ella.to/pkg/sliceutil"
)

type CustomError struct {
	Name       string
	Code       int64
	HttpStatus string
	Msg        string
}

type CustomErrors []CustomError

func (c *CustomErrors) Parse(prog *ast.Program) error {
	*c = sliceutil.Mapper(astutil.GetCustomErrors(prog), func(customError *ast.CustomError) CustomError {
		return CustomError{
			Name:       customError.Name.String(),
			Code:       customError.Code,
			HttpStatus: fmt.Sprintf("http.Status%s", ast.HttpStatusCode2String[customError.HttpStatus]),
			Msg:        customError.Msg.Value,
		}
	})

	return nil
}
