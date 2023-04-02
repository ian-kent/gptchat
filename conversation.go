package main

import (
	"fmt"
	"github.com/ian-kent/gptchat/util"
	"github.com/sashabaranov/go-openai"
	"time"
)

var systemPrompt = `You are a helpful assistant.

You enjoy conversations with the user and like asking follow up questions to gather more information.

You have commands available which you can use to help me.

You can call these commands using the slash command syntax, for example, this is how you call the help command:

` + util.TripleQuote + `
/help
` + util.TripleQuote + `

The /help command will give you a list of the commands you have available.

Commands can also include a request body, for example, this is an example of a command which takes an input:

` + util.TripleQuote + `
/example
{
    "expr": "value"
}
` + util.TripleQuote + `

Most commands also have subcommands, and this is an example of how you call a subcommand:

` + util.TripleQuote + `
/example subcommand
{
    "expr": "value"
}
` + util.TripleQuote + `

To call a command, include the command in your response. You don't need to explain the command response to me, I don't care what it is, I only care that you can use it's output to follow my instructions.`

const openingPrompt = `Hello! Please familiarise yourself with the commands you have available.

You must do this before we have a conversation.`

func intervalPrompt() string {
	return fmt.Sprintf(`The current date and time is %s.

Remember that the '/help' command will tell you what commands you have available.`, time.Now().Format("02 January 2006, 03:04pm"))
}

var conversation []openai.ChatCompletionMessage

func appendMessage(role string, message string) {
	conversation = append(conversation, openai.ChatCompletionMessage{
		Role:    role,
		Content: message,
	})
}

func resetConversation() {
	conversation = []openai.ChatCompletionMessage{}
}
