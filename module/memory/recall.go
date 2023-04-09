package memory

import (
	"context"
	"encoding/json"

	"github.com/ian-kent/gptchat/util"
	"github.com/sashabaranov/go-openai"
)

func (m *Module) Recall(input string) (string, error) {
	b, err := json.Marshal(m.memories)
	if err != nil {
		return "", err
	}

	resp, err := m.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: m.cfg.OpenAIAPIModel(),
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You are a helpful assistant.

I'll give you a list of existing memories, and a prompt which asks you to identify the memory I'm looking for.

You should review the listed memories and suggest which memories might match the request.`,
				},
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `Here are your memories in JSON format:

` + util.TripleQuote + `
` + string(b) + `
` + util.TripleQuote,
				},
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `Help me find any memories which may match this request:

` + util.TripleQuote + `
` + input + `
` + util.TripleQuote,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	response := resp.Choices[0].Message.Content
	return `You have successfully recalled this memory:

` + util.TripleQuote + `
` + response + `
` + util.TripleQuote, nil

	// TODO find a prompt which gets GPT to adjust relative time

	//	return `You have successfully recalled this memory:
	//
	//` + util.TripleQuote + `
	//` + response + `
	//` + util.TripleQuote + `
	//
	//If this memory mentions relative time (for example today, yesterday, last week, tomorrow), remember to take this into consideration when using this information to answer questions.
	//
	//For example, if the memory says "tomorrow" and the memory was stored on 25th, the memory is actually referring to 26th.`, nil
}
