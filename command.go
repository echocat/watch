package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	command = &commandHolder{}

	errNoCommand = errors.New("no command specified")
	zero         = 0

	exitCodes = app.Flag("exitCodes", "List of exit codes to leads to a next execution. Empty means always continue.").
			Short('e').
			Ints()
	successExitCodes = app.Flag("successExitCodes", "List of exit codes that indicates successful execution.").
				Short('s').
				Default("0").
				Ints()
)

func init() {
	app.Arg("command", "The command to be executed.").
		Required().
		StringsVar(&command.Arguments)
}

type commandHolder struct {
	Arguments []string

	error    error
	executed bool
}

func (instance *commandHolder) execute() {
	if len(instance.Arguments) <= 0 {
		instance.error = errNoCommand
		return
	}
	args := make([]string, len(instance.Arguments))
	copy(args, instance.Arguments)
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

	instance.error = cmd.Run()
}

func (instance commandHolder) String() string {
	buf := new(bytes.Buffer)
	for i, arg := range instance.Arguments {
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

func (instance commandHolder) Error() error {
	return instance.error
}

func (instance commandHolder) Executed() bool {
	return instance.executed
}

func (instance commandHolder) Succeeded() bool {
	if !instance.Executed() {
		return false
	}
	if instance.Error() != nil {
		return false
	}
	return intSliceContains(*successExitCodes, instance.ExitCode())
}

func (instance commandHolder) Failed() bool {
	if !instance.Executed() {
		return false
	}
	return instance.Succeeded()
}

type withExitStatus interface {
	ExitStatus() int
}

func (instance commandHolder) ExitCode() int {
	err := instance.Error()
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

func (instance commandHolder) ResultSummary() string {
	var exitCode *int
	err := instance.Error()

	if err == nil {
		exitCode = &zero
	} else if eErr, ok := err.(*exec.ExitError); ok {
		if wes, ok := eErr.Sys().(withExitStatus); ok {
			es := wes.ExitStatus()
			exitCode = &es
		}
	}

	if exitCode != nil {
		successSummary := "Failed"
		if intSliceContains(*successExitCodes, *exitCode) {
			successSummary = "Succeeded"
		}
		return fmt.Sprintf("%s with %d", successSummary, *exitCode)
	}

	return fmt.Sprintf("Failed with %v", err)
}
