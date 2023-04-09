package memory

import (
	"errors"
	"fmt"

	"github.com/ian-kent/gptchat/config"
	"github.com/ian-kent/gptchat/util"
	openai "github.com/sashabaranov/go-openai"
)

type memory struct {
	DateStored string `json:"date_stored"`
	Memory     string `json:"memory"`
}

type Module struct {
	cfg      config.Config
	client   *openai.Client
	memories []memory
}

func (m *Module) ID() string {
	return "memory"
}

func (m *Module) Load(cfg config.Config, client *openai.Client) error {
	m.cfg = cfg
	m.client = client
	return m.loadFromFile()
}

func (m *Module) UpdateConfig(cfg config.Config) {
	m.cfg = cfg
}

func (m *Module) Execute(args, body string) (string, error) {
	switch args {
	case "store":
		return m.Store(body)
	case "recall":
		return m.Recall(body)
	default:
		return "", errors.New(fmt.Sprintf("command not implemented: /memory %s", args))
	}
}

func (m *Module) Prompt() string {
	return memoryPrompt
}

const memoryPrompt = `You also have a working long term memory.

You can remember something using the '/memory store' command, or you can recall it using the '/memory recall' command.

For example, if you want to store a fact I give you, you might reply with a message such as:

` + util.TripleQuote + `
/memory store {
	"memory": "I bought cookies yesterday",
	"context": "The user was discussing what they'd eaten"
}
` + util.TripleQuote + `

You should always include useful context in the memory to help with recall later, for example any relevant parts of the conversation which might help you remember why you stored a particular memory.

Later, you may try to recall this memory. You only need to use related concepts, since of course until you call '/memory recall' you cannot possibly know the memory you've stored. For example:

` + util.TripleQuote + `
/memory recall {
	When did I buy cookies?
}
` + util.TripleQuote + `

or, another example might be:

` + util.TripleQuote + `
/memory recall {
	What did I buy yesterday?
}
` + util.TripleQuote + `

This means you can remember interesting facts about our conversation. You should always consult your long term memory using the '/memory recall' command if you think I'm asking about something you don't have information about. For example, if I ask you "What did I eat for lunch yesterday?", and you don't know the answer, you must try to retrieve the answer from your long term memory.

You can use '/memory store' and '/memory recall' at any time. You might find this useful if I ask you to remember something. The commands must be an entire message, with no conversational text, and no quotes.

If I give you information that you don't already know, you must store this information in your memory. This could be any sort of information, for example a story, information about myself, people I know, things I've done, real world events. If you're unsure, you should store it just in case.

You must not remember the current date. The current date changes and is not a useful memory.`
