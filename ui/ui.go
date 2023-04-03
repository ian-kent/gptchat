package ui

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

var (
	cUsername = color.New(color.FgRed)
	cMessage  = color.New(color.FgBlue)
	cUseful   = color.New(color.FgWhite)
	cAI       = color.New(color.FgGreen)
	cUser     = color.New(color.FgYellow)
	cError    = color.New(color.FgHiRed, color.Bold)
	cWarn     = color.New(color.FgHiYellow, color.Bold)

	cApp     = color.New(color.FgWhite)
	cAppBold = color.New(color.FgGreen, color.Bold)
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
	cError.Printf("ERROR: ")
	cUseful.Printf("%s: %v\n\n", message, err)
}

func Warn(message string) {
	cWarn.Printf("WARNING: ")
	cUseful.Printf("%s\n", message)
}

func Welcome(title, message string) {
	cAppBold.Printf("%s\n\n", title)
	cApp.Printf("%s\n\n", message)
}

func PrintChatDebug(name, message string) {
	cUseful.Printf("[DEBUG] ")
	PrintChat(name, message)
}

func PrintChat(name, message string) {
	switch name {
	case User:
		cUser.Printf("%s:\n\n", name)
		cMessage.Printf("%s\n", indent(message))
	case AI:
		cAI.Printf("%s:\n\n", name)
		cUseful.Printf("%s\n", indent(message))
	case App:
		cAppBold.Printf("%s:\n\n", name)
		cUseful.Printf("%s\n", indent(message))
	case System:
		fallthrough
	case Tool:
		fallthrough
	case API:
		fallthrough
	case Module:
		fallthrough
	default:
		cUsername.Printf("%s:\n\n", name)
		cMessage.Printf("%s\n", indent(message))
	}
}

func PromptChatInput() string {
	reader := bufio.NewReader(os.Stdin)
	cUser.Printf("USER:\n\n    ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	fmt.Println()

	return text
}

func PromptConfirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	cAppBold.Printf("%s [Y/N]: ", prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	fmt.Println()

	return strings.ToUpper(text) == "Y"
}

func PromptInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	cAppBold.Printf("%s ", prompt)
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
