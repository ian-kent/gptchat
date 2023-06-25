package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	User   = "USER"
	AI     = "AI"
	System = "SYSTEM"
	Tool   = "TOOL"
	API    = "API"
	Module = "MODULE"
	App    = "APP"
)

func Error(message string, err error) {
	theme.Error.Printf("ERROR: ")
	theme.Useful.Printf("%s: %v\n\n", message, err)
}

func Warn(message string) {
	theme.Warn.Printf("WARNING: ")
	theme.Useful.Printf("%s\n", message)
}

func Info(message string) {
	theme.Warn.Printf("INFO: ")
	theme.Useful.Printf("%s\n", message)
}

func Welcome(title, message string) {
	theme.AppBold.Printf("%s\n\n", title)
	theme.App.Printf("%s\n\n", message)
}

func PrintChatDebug(name, message string) {
	theme.Useful.Printf("[DEBUG] ")
	PrintChat(name, message)
}

func PrintChat(name, message string) {
	switch name {
	case User:
		theme.User.Printf("%s:\n\n", name)
		theme.Message.Printf("%s\n", indent(message))
	case AI:
		theme.AI.Printf("%s:\n\n", name)
		theme.Useful.Printf("%s\n", indent(message))
	case App:
		theme.AppBold.Printf("%s:\n\n", name)
		theme.Useful.Printf("%s\n", indent(message))
	case System:
		fallthrough
	case Tool:
		fallthrough
	case API:
		fallthrough
	case Module:
		fallthrough
	default:
		theme.Username.Printf("%s:\n\n", name)
		theme.Message.Printf("%s\n", indent(message))
	}
}

func PromptChatInput() string {
	reader := bufio.NewReader(os.Stdin)
	theme.User.Printf("USER:\n\n    ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	fmt.Println()

	return text
}

func PromptConfirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	theme.AppBold.Printf("%s [Y/N]: ", prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	fmt.Println()

	return strings.ToUpper(text) == "Y"
}

func PromptInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	theme.AppBold.Printf("%s ", prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func indent(input string) string {
	lines := strings.Split(string(input), "\n")
	var output string
	for _, line := range lines {
		output += "    " + line + "\n"
	}
	return output
}
