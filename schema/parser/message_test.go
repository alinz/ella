package parser_test

import (
	"testing"
)

func TestMessageParser(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `message A {
				Value: map<string, string> # hello world
			}`,
			Output: `
message A {
	Value: map<string, string>
}			
			`,
		},
		{
			Input: `message A {
				Value: map<string, string> {
					json = "hello's world"
				}
			}`,
			Output: `
message A {
	Value: map<string, string> {
		json = "hello's world"
	}
}			
			`,
		},
		{
			Input: `message A {
				Value: int64 {
					json = "hello's world"
				}
			}`,
			Output: `
message A {
	Value: int64 {
		json = "hello's world"
	}
}			
			`,
		},
		{
			Input: `message A {
				Value: int64 {
					json = 'hello world'
				}
			}`,
			Output: `
message A {
	Value: int64 {
		json = "hello world"
	}
}			
			`,
		},
		{
			Input: `message A {
				Value: int64 {
					json = "hello world"
				}
			}`,
			Output: `
message A {
	Value: int64 {
		json = "hello world"
	}
}			
			`,
		},
		{
			Input: `message A {
				Value: int64 {
					json = "value"
				}
			}`,
			Output: `
message A {
	Value: int64 {
		json = "value"
	}
}			
			`,
		},
		{
			Input: `message B {
				Value: int64
				...A
			}`,
			Output: `
message B {
	...A
	Value: int64
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
				Value: int32
			}`,
			Output: `
message B {
	Value: int32
}`,
		},
		{
			Input:  `message A {}`,
			Output: `message A {}`,
		},
	})
}
