package validator

import (
	"fmt"

	"ella.to/schema/ast"
)

func Validate(program *ast.Program) error {
	// need to make sure no duplicate names defines

	// need to sort nodes
	// 1. constants
	// 2. enums
	// 3. messages
	// 4. services

	constantsMap := make(map[string]*ast.Constant)
	enumsMap := make(map[string]*ast.Enum)
	messagesMap := make(map[string]*ast.Message)
	servicesMap := make(map[string]*ast.Service)

	{
		// create a scope since we don't need to check for duplicate names
		names := make(map[string]struct{})

		for _, node := range program.Nodes {
			switch node := node.(type) {
			case *ast.Constant:
				if _, ok := names[node.Name.Name]; ok {
					return fmt.Errorf("constant already defined: %s", node.Name.Name)
				}
				names[node.Name.Name] = struct{}{}
				constantsMap[node.Name.Name] = node

			case *ast.Enum:
				if _, ok := names[node.Name.Name]; ok {
					return fmt.Errorf("enum already defined: %s", node.Name.Name)
				}
				names[node.Name.Name] = struct{}{}
				enumsMap[node.Name.Name] = node

			case *ast.Message:
				if _, ok := names[node.Name.Name]; ok {
					return fmt.Errorf("message already defined: %s", node.Name)
				}
				names[node.Name.Name] = struct{}{}
				messagesMap[node.Name.Name] = node

			case *ast.Service:
				if _, ok := names[node.Name.Name]; ok {
					return fmt.Errorf("service already defined: %s", node.Name)
				}
				names[node.Name.Name] = struct{}{}
				servicesMap[node.Name.Name] = node
			}
		}
	}

	constants, err := validateConstants(constantsMap)
	if err != nil {
		return err
	}

	enums, err := validateEnums(enumsMap)
	if err != nil {
		return err
	}

	messages, err := validateMessages(messagesMap, enumsMap)
	if err != nil {
		return err
	}

	services, err := validateServices(servicesMap, messagesMap, enumsMap)
	if err != nil {
		return err
	}

	program.Nodes = make([]ast.Node, 0, len(constants)+len(enums)+len(messages)+len(services))

	for _, constant := range constants {
		program.Nodes = append(program.Nodes, constant)
	}

	for _, enum := range enums {
		program.Nodes = append(program.Nodes, enum)
	}

	for _, message := range messages {
		program.Nodes = append(program.Nodes, message)
	}

	for _, service := range services {
		program.Nodes = append(program.Nodes, service)
	}

	return nil
}
