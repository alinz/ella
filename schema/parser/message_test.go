package parser_test

import (
	"testing"
)

func TestMessageParser(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `message A {
				value: map<String, String> {
					json = "hello's world"
				}
			}`,
			Output: `
message A {
	value: map<String, String> {
		json = "hello's world"
	}
}			
			`,
		},
		{
			Input: `message A {
				value: Int64 {
					json = "hello's world"
				}
			}`,
			Output: `
message A {
	value: Int64 {
		json = "hello's world"
	}
}			
			`,
		},
		{
			Input: `message A {
				value: Int64 {
					json = 'hello world'
				}
			}`,
			Output: `
message A {
	value: Int64 {
		json = "hello world"
	}
}			
			`,
		},
		{
			Input: `message A {
				value: Int64 {
					json = "hello world"
				}
			}`,
			Output: `
message A {
	value: Int64 {
		json = "hello world"
	}
}			
			`,
		},
		{
			Input: `message A {
				value: Int64 {
					json = "value"
				}
			}`,
			Output: `
message A {
	value: Int64 {
		json = value
	}
}			
			`,
		},
		{
			Input: `message B {
				value: int64
				...A
			}`,
			Output: `
message B {
	...A
	value: int64
}`,
		},
		{
			Input: `message B {
				...A
			}`,
			Output: `
message B {
	...A
}`,
		},
		{
			Input: `message B {
				value: Int32
			}`,
			Output: `
message B {
	value: Int32
}`,
		},
		{
			Input:  `message A {}`,
			Output: `message A {}`,
		},
	})
}
