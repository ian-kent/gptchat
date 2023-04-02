package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output ParseResult
	}{
		{
			name:  "basic Command",
			input: "/api get /path",
			output: ParseResult{
				Chat: "",
				Commands: []ParseCommand{
					{
						Command: "/api",
						Args:    "get /path",
						Body:    "",
					},
				},
			},
		},
		{
			name:  "Command with args",
			input: "/api get /path { something }",
			output: ParseResult{
				Chat: "",
				Commands: []ParseCommand{
					{
						Command: "/api",
						Args:    "get /path",
						Body:    "{ something }",
					},
				},
			},
		},
		{
			name: "multiline Command with args",
			input: `/api get /path {
	something
}`,
			output: ParseResult{
				Chat: "",
				Commands: []ParseCommand{
					{
						Command: "/api",
						Args:    "get /path",
						Body: `{
	something
}`,
					},
				},
			},
		},
		{
			name: "multiline Command with args",
			input: `/api get /path
{
	something
}`,
			output: ParseResult{
				Chat: "",
				Commands: []ParseCommand{
					{
						Command: "/api",
						Args:    "get /path",
						Body: `{
	something
}`,
					},
				},
			},
		},
		{
			name: "chat with multiline Command with args",
			input: `This is some chat

/api get /path
{
	something
}`,
			output: ParseResult{
				Chat: "This is some chat",
				Commands: []ParseCommand{
					{
						Command: "/api",
						Args:    "get /path",
						Body: `{
	something
}`,
					},
				},
			},
		},
		{
			name: "chat with multiple multiline Command with args",
			input: `This is some chat

/api get /path
{
	something
}

This is some more chat

/api post /another-path
{
	something else
}`,
			output: ParseResult{
				Chat: "This is some chat\n\nThis is some more chat",
				Commands: []ParseCommand{
					{
						Command: "/api",
						Args:    "get /path",
						Body: `{
	something
}`,
					},
					{
						Command: "/api",
						Args:    "post /another-path",
						Body: `{
	something else
}`,
					},
				},
			},
		},
		{
			name: "multiline Command with code Body",
			input: `/plugins create my-plugin
{
	package main

	import "fmt"

	func main() {
		fmt.Println("test")
	}
}`,
			output: ParseResult{
				Chat: "",
				Commands: []ParseCommand{
					{
						Command: "/plugins",
						Args:    "create my-plugin",
						Body: `{
	package main

	import "fmt"

	func main() {
		fmt.Println("test")
	}
}`,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tokens := Lex(testCase.input)
			ast := ParseTokens(tokens)
			assert.Equal(t, testCase.output, ast)
		})
	}
}
