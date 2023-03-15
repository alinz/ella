package typescript_test

import (
	"testing"

	"ella.to/schema/ast"
	"ella.to/transform"
	"ella.to/transform/typescript"
)

func TestServices(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `
Ella = "0.0.1"

service AppleService {
	Ping() => (pings: stream int64)
	GetStatus(id: string) => (status: string)
}
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				services := getValues[*ast.Service](prog)

				return []transform.Func{
					typescript.Services(services),
				}
			},
			Output: `
// Services

interface AppleServicePingArgs {
}

interface AppleServicePingStream {
	pings: number
}

interface AppleServiceGetStatusArgs {
	id: string
}

interface AppleServiceGetStatusReturns {
	status: string
}

export interface AppleService {
	ping: (
		args: AppleServicePingArgs,
		headers?: Record<string, string>
	) => Promise<Subscription<AppleServicePingStream>>;
	getStatus: (
		args: AppleServiceGetStatusArgs,
		headers?: Record<string, string>
	) => Promise<AppleServiceGetStatusReturns>;
}

export function createAppleServiceClient(host: string): AppleService {
	return {
	ping: (
		args: AppleServicePingArgs,
		headers?: Record<string, string>
	) => Promise<Subscription<AppleServicePingStream>> {
			return callServiceStreamMethod(
				host,
				"/AppleService/Ping",
				"POST",
				args,
				headers
			);
		},
	getStatus: (
		args: AppleServiceGetStatusArgs,
		headers?: Record<string, string>
	) => Promise<AppleServiceGetStatusReturns> {
			return callServiceMethod(
				host,
				"/AppleService/GetStatus",
				"POST",
				args,
				headers
			);
		},
	};
}			
			`,
		},
		{
			Input: `
Ella = "0.0.1"

service AppleService {
	Ping()
	GetStatus(id: string) => (status: string)
}
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				services := getValues[*ast.Service](prog)

				return []transform.Func{
					typescript.Services(services),
				}
			},
			Output: `
// Services

interface AppleServicePingArgs {
}

interface AppleServicePingReturns {
}

interface AppleServiceGetStatusArgs {
	id: string
}

interface AppleServiceGetStatusReturns {
	status: string
}

export interface AppleService {
	ping: (
		args: AppleServicePingArgs,
		headers?: Record<string, string>
	) => Promise<AppleServicePingReturns>;
	getStatus: (
		args: AppleServiceGetStatusArgs,
		headers?: Record<string, string>
	) => Promise<AppleServiceGetStatusReturns>;
}

export function createAppleServiceClient(host: string): AppleService {
	return {
	ping: (
		args: AppleServicePingArgs,
		headers?: Record<string, string>
	) => Promise<AppleServicePingReturns> {
			return callServiceMethod(
				host,
				"/AppleService/Ping",
				"POST",
				args,
				headers
			);
		},
	getStatus: (
		args: AppleServiceGetStatusArgs,
		headers?: Record<string, string>
	) => Promise<AppleServiceGetStatusReturns> {
			return callServiceMethod(
				host,
				"/AppleService/GetStatus",
				"POST",
				args,
				headers
			);
		},
	};
}
			`,
		},
	})
}
