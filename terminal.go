package main

import (
	"github.com/mattn/go-colorable"
	"io"
	"os"
)

var (
	colored = app.Flag("colored", `"auto" - Will decide the best option
| "always" - Will always use colors
| "never" - Will never use colors`).
		Short('c').
		Default("auto").
		Enum("always", "auto", "never")
	baseTerminalOutput io.Writer = os.Stdout
	terminalOutput     io.Writer = os.Stderr
)

func initTerminal() {
	baseTerminalOutput = colorable.NewColorable(os.Stdout)
	terminalOutput = colorable.NewColorable(os.Stderr)
}

func resetTerminal() {
	switch *colored {
	case "auto", "always":
		const clear = "\u001b[2J\u001b[1;1H"
		mustFprint(baseTerminalOutput, clear)
		mustFprint(terminalOutput, clear)
	}
}

func printHighlightedToTerminal(what string) {
	switch *colored {
	case "auto", "always":
		mustFprint(terminalOutput, "\u001b[47m\u001b[30m", what, "\u001b[0m")
	default:
		mustFprint(terminalOutput, what)
	}

}
