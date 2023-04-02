package main

import (
	"fmt"
	"github.com/ian-kent/gptchat/parser"
	"strings"
)

func main() {
	input := `/plugins create my-plugin
{
	package main

	import "fmt"

	func main() {
		fmt.Println("test")
	}
}`

	tokens := parser.Lex(input)

	fmt.Println("Tokens:")
	for _, token := range tokens {
		fmt.Printf("    %20s => %s\n", token.Typ, token.Val)
	}

	fmt.Println()

	result := parser.ParseTokens(tokens)
	fmt.Println("Result:")
	fmt.Println("    Chat:")
	fmt.Println(indent(result.Chat, "        "))
	fmt.Println("    Commands:")
	for _, command := range result.Commands {
		fmt.Printf("        - Command: %s\n", command.Command)
		fmt.Printf("        - Args: %s\n", command.Args)
		fmt.Printf("        - Body:\n")
		fmt.Println(indent(command.Body, "              "))
	}
}

func indent(input string, prefix string) string {
	lines := strings.Split(string(input), "\n")
	var output string
	for _, line := range lines {
		output += prefix + line + "\n"
	}
	return output
}
