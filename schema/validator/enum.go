package validator

import (
	"fmt"
	"sort"

	"ella.to/schema/ast"
)

type counter struct {
	size            int
	unsigned        bool
	unsignedCounter uint64
	signedCounter   int64
}

func (c *counter) increment(ptr *ast.Constant) {
	if c.unsigned {
		ptr.Value = &ast.ValueUint{
			Size:    c.size,
			Content: c.unsignedCounter,
		}
		c.unsignedCounter++
	} else {
		ptr.Value = &ast.ValueInt{
			Size:    c.size,
			Content: c.signedCounter,
		}
		c.signedCounter++
	}
}

func (c *counter) set(val any) {
	if c.unsigned {
		c.unsignedCounter = val.(uint64)
		c.unsignedCounter++
	} else {
		c.signedCounter = val.(int64)
		c.signedCounter++
	}
}

func newCounter(size int, unsigned bool) *counter {
	return &counter{
		size:     size,
		unsigned: unsigned,
	}
}

func validateEnum(enum *ast.Enum) error {
	var counter *counter

	switch v := enum.Type.(type) {
	case *ast.TypeInt:
		counter = newCounter(v.Size, false)
	case *ast.TypeUint:
		counter = newCounter(v.Size, true)
	}

	for _, constant := range enum.Constants {
		switch cv := constant.Value.(type) {
		case *ast.ValueInt:
			counter.set(cv.Content)
		case *ast.ValueUint:
			counter.set(cv.Content)
		case nil:
			counter.increment(constant)
		default:
			return fmt.Errorf("enum constant %s must be integer", constant.Name.Name)
		}
	}

	duplicateCheck := make(map[string]struct{})

	for _, constant := range enum.Constants {
		key := constant.Value.TokenLiteral()
		if _, ok := duplicateCheck[key]; ok {
			return fmt.Errorf("enum constant %s is duplicate name or value", constant.Name.Name)
		}
		duplicateCheck[key] = struct{}{}
	}

	return nil
}

func validateEnums(enumsMap map[string]*ast.Enum) ([]*ast.Enum, error) {
	enums := make([]*ast.Enum, 0, len(enumsMap))

	for _, enum := range enumsMap {
		if err := validateEnum(enum); err != nil {
			return nil, err
		}
		enums = append(enums, enum)
	}

	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name.Name < enums[j].Name.Name
	})

	return enums, nil
}
