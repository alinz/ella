package typescript

import (
	"ella.to/schema/ast"
	"ella.to/transform"
)

func ServiceMethodArgs(serviceName string, method *ast.Method) transform.Func {
	return func(out transform.Writer) error {
		out.
			Str(`interface `).Pascal(serviceName + method.Name.Name).Str(`Args {`).
			Lines(1)

		for _, arg := range method.Args {
			typ, err := parseType(arg.Type)
			if err != nil {
				return err
			}
			out.
				Tabs(1).
				Snake(arg.Name.Name).Str(`: `).Str(typ).
				Lines(1)
		}

		out.Str(`}`)

		return nil
	}
}

func ServiceMethodReturns(serviceName string, method *ast.Method) transform.Func {
	return func(out transform.Writer) error {
		hasStream := containsStream(method)

		// the follwing operation created 2 different lines based on hasStream
		// interface Service1Method1Stream {
		// interface Service1Method1Returns {
		out.
			Str("interface ").
			Pascal(serviceName+method.Name.Name).
			StrCond(hasStream, "Stream").
			StrCond(!hasStream, "Returns").
			Str(" {").
			Lines(1)

		for _, ret := range method.Returns {
			typ, err := parseType(ret.Type)
			if err != nil {
				return err
			}
			out.
				Tabs(1).
				Snake(ret.Name.Name).Str(`: `).Str(typ).
				Lines(1)
		}

		out.Str(`}`)

		return nil
	}
}

func ServiceInterface(service *ast.Service) transform.Func {
	return func(out transform.Writer) error {
		serviceName := service.Name.Name

		out.
			Str(`export interface `).Pascal(serviceName).Str(` {`).
			Lines(1)

		for _, method := range service.Methods {
			hasStream := containsStream(method)

			out.
				Tabs(1).
				Camel(method.Name.Name).Str(`: (`).
				Lines(1)

			out.
				Tabs(2).Str(`args: `).Pascal(serviceName + method.Name.Name).Str(`Args,`).Lines(1).
				Tabs(2).Str(`headers?: Record<string, string>`).Lines(1)

			// the follwing operation created 2 different lines based on hasStream
			// ) => Promise<Subscription<Service1Method1Stream>>;
			// ) => Promise<Service1Method1Returns>;
			out.
				Tabs(1).
				Str(`) => Promise<`).
				StrCond(hasStream, "Subscription<").
				Pascal(serviceName+method.Name.Name).
				StrCond(hasStream, "Stream>").
				StrCond(!hasStream, "Returns").
				Str(">;")

			out.Lines(1)
		}

		out.Str(`}`)

		return nil
	}
}

func ServiceClient(service *ast.Service) transform.Func {
	return func(out transform.Writer) error {
		serviceName := service.Name.Name

		out.
			Str(`export function create`).
			Pascal(serviceName).
			Str(`Client(host: string): `).
			Pascal(serviceName).
			Str(` {`).Lines(1)

		out.Tabs(1).Str(`return {`).Lines(1)

		for _, method := range service.Methods {
			methodName := method.Name.Name
			hasStream := containsStream(method)

			out.
				Tabs(1).
				Camel(methodName).Str(`: (`).
				Lines(1)

			out.
				Tabs(2).Str(`args: `).Pascal(serviceName + methodName).Str(`Args,`).Lines(1).
				Tabs(2).Str(`headers?: Record<string, string>`).Lines(1)

			out.
				Tabs(1).
				Str(`) => Promise<`).
				StrCond(hasStream, "Subscription<").
				Pascal(serviceName+method.Name.Name).
				StrCond(hasStream, "Stream>").
				StrCond(!hasStream, "Returns").
				Str("> {").Lines(1)

			out.
				Tabs(3).Str(`return callService`).StrCond(hasStream, "Stream").Str("Method(").Lines(1)

			var httpMethod string
			httpMethodOpt := getConstByKey(method.Options, "http.method")
			if httpMethodOpt != nil {
				httpMethod = getConstValueAsString(httpMethodOpt)
			}
			if httpMethod == "" {
				httpMethod = "POST"
			}

			out.Tabs(4).Str(`host,`).Lines(1)
			out.Tabs(4).Str(`"/`).Pascal(serviceName).Str(`/`).Pascal(methodName).Str(`",`).Lines(1)
			out.Tabs(4).Str(`"`).Str(httpMethod).Str(`",`).Lines(1)
			out.Tabs(4).Str(`args,`).Lines(1)
			out.Tabs(4).Str(`headers`).Lines(1)

			out.Tabs(3).Str(`);`).Lines(1)

			out.Tabs(2).Str(`},`).Lines(1)
		}

		out.Tabs(1).Str(`};`).Lines(1)
		out.Str(`}`)

		return nil
	}
}

func Service(service *ast.Service) transform.Func {
	return func(out transform.Writer) error {
		serviceName := service.Name.Name

		// define arg interface
		for _, method := range service.Methods {
			if err := ServiceMethodArgs(serviceName, method)(out); err != nil {
				return err
			}

			out.Lines(2)

			if err := ServiceMethodReturns(serviceName, method)(out); err != nil {
				return err
			}

			out.Lines(2)
		}

		if err := ServiceInterface(service)(out); err != nil {
			return err
		}

		out.Lines(2)

		// define service client implementation
		if err := ServiceClient(service)(out); err != nil {
			return err
		}

		return nil
	}
}

func Services(services []*ast.Service) transform.Func {
	return func(out transform.Writer) error {
		out.Str("// Services").Lines(2)

		for _, service := range services {
			if err := Service(service)(out); err != nil {
				return err
			}
			out.Lines(1)
		}

		return nil
	}
}
