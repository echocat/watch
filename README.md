# watch

[![Build Status](https://travis-ci.org/echocat/watch.svg?branch=master)](https://travis-ci.org/echocat/watch)

> Like the unix one but works cross-platform without magic.

## Install

Either download the latest binary executable from [releases page](https://github.com/echocat/watch/releases/latest) or install master using golang directly:

```
$ go install github.com/echocat/watch@latest
```

## Usage

```
usage: watch [<flags>] <command>...

Like the unix one but works cross-platform without magic.

Flags:
      --help                     Show context-sensitive help (also try
                                 --help-long and --help-man).
  -e, --exitCodes=EXITCODES ...  List of exit codes to leads to a next
                                 execution. Empty means always continue.
  -s, --successExitCodes=0 ...   List of exit codes that indicates successful
                                 execution.
  -n, --interval=5s              Execute every n duration.
  -r, --resetTerminal            If enabled before each execution the terminal
                                 will be reset.
      --timeFormat="2006-01-02 15:04:05"
                                 How to format the time. See:
                                 https://golang.org/pkg/time/
  -c, --colored=auto             "auto" - Will decide the best option | "always"
                                 - Will always use colors | "never" - Will never
                                 use colors
      --version                  Will print information about this version.
  -h, --header=[{{.Now}}] Execute [{{.Command}}] every {{.Interval}}

                                 Will print a header what will be executed. If
                                 empty no header will be displayed. See:
                                 https://golang.org/pkg/text/template/
  -t, --footer=[{{.Now}}] {{.Command.ResultSummary}}

                                 Will print a footer what was executed. If empty
                                 no footer will be displayed. See:
                                 https://golang.org/pkg/text/template/

Args:
  <command>  The command to be executed.
```

## Examples

```
$ watch kubectl -n foobar get pods
```
