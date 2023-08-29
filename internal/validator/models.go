package validator

import (
	"fmt"
	"slices"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/astutil"
)

// Validates the messages of the program from the following aspects:
// - check if cycles exists in the message extends
// - adding fields from extended messages
//   - make sure all the extends fields are unique
//
// - type of the fields must be a primitive type, base, enum or message
// - name of the message has to be PascalCase
// - name of the fields has to be PascalCase
// - name of the fields has to be unique per message
//   - if true, type must be the same and this indicates options needs to be overwritten
func validateModels(prog *ast.Program) error {
	return runValidators(
		prog,
		checkModelCycles,
		mergeExtendFields,
	)
}

func checkModelCycles(prog *ast.Program) error {
	messages := astutil.GetModels(prog)
	dependencyGraph := make(map[string][]string)

	for _, message := range messages {
		dependencyGraph[message.Name.String()] = make([]string, 0, len(message.Extends))
		for _, extend := range message.Extends {
			dependencyGraph[message.Name.String()] = append(dependencyGraph[message.Name.String()], extend.String())
		}
	}

	cycleName, hasCycles := func(graph map[string][]string) (string, bool) {
		visited := make(map[string]bool)
		recStack := make(map[string]bool)

		var hasCycle bool

		var dfsVisit func(name string)

		dfsVisit = func(name string) {
			visited[name] = true
			recStack[name] = true

			for _, dep := range graph[name] {
				if !visited[dep] {
					dfsVisit(dep)
				} else if recStack[dep] {
					hasCycle = true
					return
				}
			}

			recStack[name] = false
		}

		for name := range graph {
			if !visited[name] {
				dfsVisit(name)
				if hasCycle {
					return name, true
				}
			}
		}

		return "", false
	}(dependencyGraph)

	if hasCycles {
		return fmt.Errorf("message %s is part of a cycle", cycleName)
	}

	return nil
}

func mergeExtendFields(prog *ast.Program) error {
	messages := astutil.GetModels(prog)
	messagesMap := astutil.CreateModelTypeMap(messages)
	isValidType := astutil.CreateIsValidType(prog)
	constantsMap := astutil.CreateConstsMap(prog)

	for _, message := range messages {
		// check if all the extends are uniques
		extends := make(map[string]struct{})

		for _, extend := range message.Extends {
			if _, ok := extends[extend.String()]; ok {
				return fmt.Errorf("message %s is extending %s multiple times", message.Name, extend)
			}
			extends[extend.String()] = struct{}{}

			baseModel, ok := messagesMap[extend.String()]
			if !ok {
				return fmt.Errorf("message %s is extending unknown message %s", message.Name, extend)
			}

			if err := mergeFields(message, baseModel, isValidType, constantsMap); err != nil {
				return err
			}
		}
	}

	return nil
}

func mergeFields(target *ast.Model, base *ast.Model, isValidType func(typ ast.Type) bool, constantsMap map[string]*ast.Const) error {
	// append all the base fields at the beginning of the target fields
	target.Fields = append(base.Fields, target.Fields...)

	// check the fields type
	for _, field := range target.Fields {
		if !isValidType(field.Type) {
			return fmt.Errorf("message %s has a field %s with an invalid type %s", target.Name, field.Name, field.Type)
		}
	}

	fieldsMap := make(map[string]*ast.Field)
	for _, field := range target.Fields {
		baseFiled, ok := fieldsMap[field.Name.String()]
		if !ok {
			err := prepareFieldOptions(field, constantsMap)
			if err != nil {
				return err
			}
			fieldsMap[field.Name.String()] = field
			continue
		}

		if baseFiled.Type.String() != field.Type.String() {
			return fmt.Errorf("message %s has a field %s with a different type %s", target.Name, field.Name, field.Type)
		}

		err := mergeFieldOptions(baseFiled, field, constantsMap)
		if err != nil {
			return err
		}
	}

	target.Fields = make([]*ast.Field, 0, len(fieldsMap))
	for _, field := range fieldsMap {
		target.Fields = append(target.Fields, field)
	}

	slices.SortFunc(target.Fields, func(a, b *ast.Field) int {
		return strings.Compare(a.Name.String(), b.Name.String())
	})

	// check all the fields are unique
	// if not, check if the type is the same
	// if not, return an error

	return nil
}

func prepareFieldOptions(field *ast.Field, constantsMap map[string]*ast.Const) error {
	for _, option := range field.Options {
		value, err := getValue(option.Value, constantsMap)
		if err != nil {
			return err
		}

		option.Value = value
	}

	return nil
}

func mergeFieldOptions(target, base *ast.Field, constantsMap map[string]*ast.Const) error {
	target.Options = append(base.Options, target.Options...)

	optionsMap := make(map[string]*ast.Option)
	for _, option := range target.Options {
		baseOpt, ok := optionsMap[option.Name.String()]
		if !ok {
			optionsMap[option.Name.String()] = option
			continue
		}

		err := mergeOptionValue(baseOpt, option, constantsMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeOptionValue(target, base *ast.Option, constantsMap map[string]*ast.Const) error {
	targetValue, err := getValue(target.Value, constantsMap)
	if err != nil {
		return err
	}
	baseValue, err := getValue(base.Value, constantsMap)
	if err != nil {
		return err
	}

	targetType := astutil.GetValueType(targetValue)
	baseType := astutil.GetValueType(baseValue)
	if targetType != baseType {
		return fmt.Errorf("option %s has a different type %s", target.Name, targetType)
	}

	return nil
}

// getValueRef returns the reference of the value, but it does check whether the value
// refer to the right constants
func getValueRef(value ast.Value, constantsMap map[string]*ast.Const) (ast.Value, error) {
	switch v := value.(type) {
	case *ast.ValueVariable:
		_, err := getConstValue(v.TokenLiteral(), constantsMap)
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		return value, nil
	}
}

// getValue evaluate the value and return the value
func getValue(value ast.Value, constantsMap map[string]*ast.Const) (ast.Value, error) {
	switch v := value.(type) {
	case *ast.ValueVariable:
		return getConstValue(v.TokenLiteral(), constantsMap)
	default:
		return value, nil
	}
}

func getConstValue(name string, constantsMap map[string]*ast.Const) (ast.Value, error) {
	constant, ok := constantsMap[name]
	if !ok {
		return nil, fmt.Errorf("unknown constant %s", name)
	}

	point, ok := constant.Value.(*ast.ValueVariable)
	if ok {
		return getConstValue(point.TokenLiteral(), constantsMap)
	}

	return constant.Value, nil
}
