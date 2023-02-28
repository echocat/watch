package watch

import (
	"github.com/alecthomas/kingpin"
	"github.com/mattn/go-colorable"
	"io"
	"os"
)

func newTerminal() *terminal {
	return &terminal{
		baseTerminalOutput: os.Stdout,
		terminalOutput:     os.Stderr,
	}
}

type terminal struct {
	colored            string
	baseTerminalOutput io.Writer
	terminalOutput     io.Writer
}

func (this *terminal) ConfigureCli(cli *kingpin.Application) {
	cli.Flag("colored", `"auto" - Will decide the best option
| "always" - Will always use colors
| "never" - Will never use colors`).
		Short('c').
		Default("auto").
		EnumVar(&this.colored, "always", "auto", "never")
}

func (this *terminal) init() {
	this.baseTerminalOutput = colorable.NewColorable(os.Stdout)
	this.terminalOutput = colorable.NewColorable(os.Stderr)
}

func (this *terminal) reset() {
	switch this.colored {
	case "auto", "always":
		const clear = "\u001b[1;1H\u001b[2J\u001b[1;1H"
		mustFprint(this.baseTerminalOutput, clear)
		mustFprint(this.terminalOutput, clear)
	}
}

func (this *terminal) printHighlightedTo(what string) {
	switch this.colored {
	case "auto", "always":
		mustFprint(this.terminalOutput, "\u001b[47m\u001b[30m", what, "\u001b[0m")
	default:
		mustFprint(this.terminalOutput, what)
	}

}
