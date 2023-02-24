package parser_test

import "testing"

func TestServiceParser(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `service A {
				Ping2() => (pong: String) {
					foo = 1
				}
			}`,
			Output: `
service A {
	Ping2() => (pong: String) {
		foo = 1
	}
}	
			`,
		},
		{
			Input: `service A {
				Ping2() {
					foo = 1
				}
			}`,
			Output: `
service A {
	Ping2() {
		foo = 1
	}
}	
			`,
		},
		{
			Input: `service A {
				Ping(timeout: Timestamp, b: Int64) => (pong: stream String)
				Ping2()
			}`,
			Output: `
service A {
	Ping(timeout: Timestamp, b: Int64) => (pong: stream String)
	Ping2()
}	
			`,
		},
		{
			Input: `service A {
				Ping(timeout: Timestamp, b: Int64) => (pong: stream String)
			}`,
			Output: `
service A {
	Ping(timeout: Timestamp, b: Int64) => (pong: stream String)
}
			`,
		},
		{
			Input: `service A {
				Ping(timeout: Timestamp) => (pong: stream String)
			}`,
			Output: `
service A {
	Ping(timeout: Timestamp) => (pong: stream String)
}
			`,
		},
		{
			Input: `service A {
				Ping(timeout: Timestamp) => (pong: String)
			}`,
			Output: `
service A {
	Ping(timeout: Timestamp) => (pong: String)
}
			`,
		},
		{
			Input: `service A {
				Ping() => (pong: String)
			}`,
			Output: `
service A {
	Ping() => (pong: String)
}
			`,
		},
		{
			Input: `service A {
				Ping()
			}`,
			Output: `
service A {
	Ping()
}
			`,
		},
		{
			Input:  `service A {}`,
			Output: `service A {}`,
		},
	})
}
