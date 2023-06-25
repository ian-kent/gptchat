package ui

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

var theme Theme

func init() {
	if strings.ToUpper(os.Getenv("GPTCHAT_THEME")) == "DARK" {
		theme = DarkTheme
	} else {
		theme = LightTheme
	}
}

type Theme struct {
	Username *color.Color
	Message  *color.Color
	Useful   *color.Color
	AI       *color.Color
	User     *color.Color
	Info     *color.Color
	Error    *color.Color
	Warn     *color.Color
	App      *color.Color
	AppBold  *color.Color
}

var LightTheme = Theme{
	Username: color.New(color.FgRed),
	Message:  color.New(color.FgBlue),
	Useful:   color.New(color.FgWhite),
	AI:       color.New(color.FgGreen),
	User:     color.New(color.FgYellow),
	Info:     color.New(color.FgWhite, color.Bold),
	Error:    color.New(color.FgHiRed, color.Bold),
	Warn:     color.New(color.FgHiYellow, color.Bold),
	App:      color.New(color.FgWhite),
	AppBold:  color.New(color.FgGreen, color.Bold),
}

var DarkTheme = Theme{
	Username: color.New(color.FgRed),
	Message:  color.New(color.FgBlue),
	Useful:   color.New(color.FgBlack),
	AI:       color.New(color.FgGreen),
	User:     color.New(color.FgMagenta),
	Info:     color.New(color.FgBlack, color.Bold),
	Error:    color.New(color.FgHiRed, color.Bold),
	Warn:     color.New(color.FgHiMagenta, color.Bold),
	App:      color.New(color.FgBlack),
	AppBold:  color.New(color.FgGreen, color.Bold),
}
