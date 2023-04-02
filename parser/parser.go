package parser

import "strings"

type ParseResult struct {
	Chat     string
	Commands []ParseCommand
}

type ParseCommand struct {
	Command string
	Args    string
	Body    string
}

func (p ParseCommand) String() string {
	output := p.Command
	if p.Args != "" {
		output += " " + p.Args
	}
	if p.Body != "" {
		output += "\n" + p.Body
	}
	return output
}

func Parse(input string) ParseResult {
	tokens := Lex(input)
	return ParseTokens(tokens)
}

type TokenType string

const (
	Plaintext TokenType = "Plaintext"
	Newline             = "Newline"
	Command             = "Command"
	Body                = "Body"
)

type Token struct {
	Typ TokenType
	Val string
}

func ParseTokens(tokens []Token) ParseResult {
	result := ParseResult{}
	var activeCommand *ParseCommand
	var commands []*ParseCommand
	var isInCommandContext bool

	for _, token := range tokens {
		switch token.Typ {
		case Plaintext:
			if activeCommand == nil {
				result.Chat += token.Val
				isInCommandContext = false
				continue
			}
			if activeCommand.Args != "" {
				// ParseTokens error
				panic("ParseTokens error: Command already has args")
			}
			activeCommand.Args = strings.TrimSpace(token.Val)
		case Newline:
			if activeCommand != nil {
				activeCommand = nil
				continue
			}
			if strings.HasSuffix(result.Chat, "\n\n") {
				// we don't append more than two consecutive newlines
				continue
			}
			result.Chat += token.Val
		case Command:
			activeCommand = &ParseCommand{Command: token.Val}
			commands = append(commands, activeCommand)
			isInCommandContext = true
		case Body:
			if activeCommand != nil {
				if activeCommand.Body != "" {
					// ParseTokens error
					panic("ParseTokens error: Command already has Body")
				}
				activeCommand.Body = token.Val
				continue
			}
			if isInCommandContext {
				lastCommand := commands[len(commands)-1]
				if lastCommand.Body != "" {
					// ParseTokens error
					panic("ParseTokens error: Command already has Body")
				}
				lastCommand.Body = token.Val
				continue
			}

			result.Chat += token.Val
		}
	}

	result.Chat = strings.TrimSpace(result.Chat)
	for _, command := range commands {
		result.Commands = append(result.Commands, *command)
	}
	return result
}

func Lex(input string) []Token {
	var tokens []Token
	var currentToken *Token
	var nesting int

	for i, c := range input {
		switch c {
		case '/':
			// is this the start of a new Command?
			if i == 0 || input[i-1] == '\n' {
				if currentToken != nil {
					tokens = append(tokens, *currentToken)
				}
				currentToken = &Token{Typ: Command, Val: "/"}
				continue
			}
			// if we have a Token, append to it
			if currentToken != nil {
				currentToken.Val += string(c)
				continue
			}
			// otherwise we can assume this is plain text
			currentToken = &Token{Typ: Plaintext, Val: "/"}
		case ' ':
			// a space signifies the end of the Command
			if currentToken != nil && currentToken.Typ == Command {
				tokens = append(tokens, *currentToken)
				currentToken = nil
				continue
			}
			// if we have a Token, append to it
			if currentToken != nil {
				currentToken.Val += string(c)
				continue
			}
			// otherwise we can assume this is plain text
			currentToken = &Token{Typ: Plaintext, Val: " "}
		case '{':
			// if it's not already a Body, we'll store the current Token
			if currentToken != nil && currentToken.Typ != Body {
				tokens = append(tokens, *currentToken)
			}
			// If we already have a body, we'll add to it
			if currentToken != nil && currentToken.Typ == Body {
				nesting++
				currentToken.Val += string(c)
				continue
			}
			// Otherwise we'll start a new body
			currentToken = &Token{Typ: Body, Val: "{"}
			nesting++
		case '}':
			// if we're already in a Body, the Body ends
			if currentToken != nil && currentToken.Typ == Body {
				nesting--
				currentToken.Val += "}"

				if nesting == 0 {
					tokens = append(tokens, *currentToken)
					currentToken = nil
				}
				continue
			}
			// if we have a Token, append to it
			if currentToken != nil {
				currentToken.Val += string(c)
				continue
			}
			// otherwise we can assume this is plain text
			currentToken = &Token{Typ: Plaintext, Val: " "}
		case '\n':
			// if we have Plaintext or a Body, we'll append to it
			if currentToken != nil && (currentToken.Typ == Body) {
				currentToken.Val += "\n"
				continue
			}

			// otherwise we always end the current Token at a new line
			if currentToken != nil {
				tokens = append(tokens, *currentToken)
				currentToken = nil
			}

			// and we store a new line Token
			tokens = append(tokens, Token{Typ: Newline, Val: "\n"})
		default:
			if currentToken == nil {
				currentToken = &Token{Typ: Plaintext}
			}
			currentToken.Val += string(c)
		}
	}
	if currentToken != nil {
		tokens = append(tokens, *currentToken)
	}

	return tokens
}
