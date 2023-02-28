package watch

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"time"
)

func NewWatch(version, revision, built string) (*Watch, error) {
	vv, err := NewVersion(version, revision, built)
	if err != nil {
		return nil, err
	}

	result := &Watch{
		Interval:                   5 * time.Second,
		ResetTerminalBeforeEachRun: true,

		Version:  *vv,
		command:  *newCommand(),
		terminal: *newTerminal(),
	}

	result.rendering = *newRendering(
		func() *command {
			return &result.command
		},
		func() time.Duration {
			return result.Interval
		},
	)

	return result, nil
}

type Watch struct {
	Interval                   time.Duration
	ResetTerminalBeforeEachRun bool

	Version   Version
	command   command
	terminal  terminal
	rendering rendering
}

func (this *Watch) ConfigureCli(cli *kingpin.Application) {
	cli.Flag("interval", "Execute every n duration.").
		Short('n').
		Default(this.Interval.String()).
		DurationVar(&this.Interval)
	cli.Flag("resetTerminal", "If enabled before each execution the terminal will be reset.").
		Short('r').
		Default(fmt.Sprint(this.ResetTerminalBeforeEachRun)).
		BoolVar(&this.ResetTerminalBeforeEachRun)
	cli.Action(func(*kingpin.ParseContext) error {
		this.Execute()
		return nil
	})

	this.Version.ConfigureCli(cli)
	this.command.ConfigureCli(cli)
	this.terminal.ConfigureCli(cli)
	this.rendering.ConfigureCli(cli)
}

func (this *Watch) Execute() {
	this.terminal.init()

	for {
		this.Run()
		time.Sleep(this.Interval)
	}
}

func (this *Watch) Run() {
	if this.ResetTerminalBeforeEachRun {
		this.terminal.reset()
	}
	this.printHeader()
	this.command.execute()
	this.printFooter()
	this.command.exitIfRequired()
}

func (this *Watch) printHeader() {
	this.terminal.printHighlightedTo(this.rendering.renderHeader())
}

func (this *Watch) printFooter() {
	this.terminal.printHighlightedTo(this.rendering.renderFooter())
}
