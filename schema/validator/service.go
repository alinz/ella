package validator

import (
	"sort"

	"github.com/alinz/rpc.go/schema/ast"
)

func validateService(service *ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	return nil
}

func validateServices(servicesMap map[string]*ast.Service, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) ([]*ast.Service, error) {
	services := make([]*ast.Service, 0, len(servicesMap))

	for _, service := range servicesMap {
		if err := validateService(service, messagesMap, enumsMap); err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].Name.Name < services[j].Name.Name
	})

	return services, nil
}
