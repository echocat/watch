package watch

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/alecthomas/kingpin"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	errNoCommand = errors.New("no command specified")
)

func newCommand() *command {
	return &command{}
}

type command struct {
	Arguments []string

	error    error
	executed bool

	exitCodes        []int
	successExitCodes []int
}

func (this *command) ConfigureCli(cli *kingpin.Application) {
	cli.Arg("command", "The command to be executed.").
		Required().
		StringsVar(&this.Arguments)
	cli.Flag("exitCodes", "List of exit codes to leads to a next execution. Empty means always continue.").
		Short('e').
		IntsVar(&this.exitCodes)
	cli.Flag("successExitCodes", "List of exit codes that indicates successful execution.").
		Short('s').
		Default("0").
		IntsVar(&this.successExitCodes)
}

func (this *command) execute() {
	if len(this.Arguments) <= 0 {
		this.error = errNoCommand
		return
	}
	args := make([]string, len(this.Arguments))
	copy(args, this.Arguments)
	cmd := &exec.Cmd{
		Path:   args[0],
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if filepath.Base(cmd.Path) == cmd.Path {
		if lp, err := exec.LookPath(cmd.Path); err != nil {
			fatal("cannot lookup '%s': %v", cmd.Path, err)
		} else {
			cmd.Path = lp
			cmd.Args[0] = lp
		}
	}

	this.error = cmd.Run()
}

func (this command) exitIfRequired() {
	if len(this.exitCodes) == 0 {
		return
	}
	exitCode := this.ExitCode()
	if !intSliceContains(this.exitCodes, exitCode) {
		if exitCode < 0 {
			exitCode = exitCode * -1
		}
		os.Exit(exitCode)
	}
}

func (this command) String() string {
	buf := new(bytes.Buffer)
	for i, arg := range this.Arguments {
		if i > 0 {
			buf.WriteString(" ")
		}
		if strings.ContainsRune(arg, '\t') ||
			strings.ContainsRune(arg, '\n') ||
			strings.ContainsRune(arg, ' ') ||
			strings.ContainsRune(arg, '\xFF') ||
			strings.ContainsRune(arg, '\u0100') ||
			strings.ContainsRune(arg, '"') ||
			strings.ContainsRune(arg, '\\') {
			buf.WriteString(strconv.Quote(arg))
		} else {
			buf.WriteString(arg)
		}
	}
	return buf.String()
}

func (this command) Error() error {
	return this.error
}

func (this command) Executed() bool {
	return this.executed
}

func (this command) Succeeded() bool {
	if !this.Executed() {
		return false
	}
	if this.Error() != nil {
		return false
	}
	return intSliceContains(this.successExitCodes, this.ExitCode())
}

func (this command) Failed() bool {
	if !this.Executed() {
		return false
	}
	return this.Succeeded()
}

type withExitStatus interface {
	ExitStatus() int
}

func (this command) ExitCode() int {
	err := this.Error()
	if err == nil {
		return 0
	}
	if eErr, ok := err.(*exec.ExitError); ok {
		if wes, ok := eErr.Sys().(withExitStatus); ok {
			return wes.ExitStatus()
		}
	}
	return -127
}

func (this command) ResultSummary() string {
	var exitCode *int
	err := this.Error()

	if err == nil {
		var zero int
		exitCode = &zero
	} else if eErr, ok := err.(*exec.ExitError); ok {
		if wes, ok := eErr.Sys().(withExitStatus); ok {
			es := wes.ExitStatus()
			exitCode = &es
		}
	}

	if exitCode != nil {
		successSummary := "Failed"
		if intSliceContains(this.successExitCodes, *exitCode) {
			successSummary = "Succeeded"
		}
		return fmt.Sprintf("%s with %d", successSummary, *exitCode)
	}

	return fmt.Sprintf("Failed with %v", err)
}
