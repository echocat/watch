package main

import (
	"github.com/alecthomas/kingpin"
	"os"
	"time"
)

var (
	app = kingpin.New("watch", "Like the unix one but works cross-platform without magic.").
		Interspersed(false)

	interval = app.Flag("interval", "Execute every n duration.").
			Short('n').
			Default("5s").
			Duration()
	resetTerminalBeforeEachRun = app.Flag("resetTerminal", "If enabled before each execution the terminal will be reset.").
					Short('r').
					Default("true").
					Bool()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	initTerminal()

	for {
		run()
		time.Sleep(*interval)
	}
}

func run() {
	if *resetTerminalBeforeEachRun {
		resetTerminal()
	}
	printHeader()
	command.execute()
	printFooter()
	exitIfRequired()
}
func exitIfRequired() {
	if len(*exitCodes) == 0 {
		return
	}
	exitCode := command.ExitCode()
	if !intSliceContains(*exitCodes, exitCode) {
		if exitCode < 0 {
			exitCode = exitCode * -1
		}
		os.Exit(exitCode)
	}
}
