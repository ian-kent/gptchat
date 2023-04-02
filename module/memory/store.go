package memory

import (
	"github.com/ian-kent/gptchat/util"
	"time"
)

func (m *Module) Store(input string) (string, error) {
	err := m.appendMemory(memory{
		DateStored: time.Now().Format("02 January 2006, 03:04pm"),
		Memory:     input,
	})
	if err != nil {
		return "", err
	}

	return `You have successfully stored this memory:

` + util.TripleQuote + `
` + input + `
` + util.TripleQuote, nil
}
