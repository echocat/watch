package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/echocat/watch/pkg/core"
)

var (
	version  = "development"
	revision = ""
	built    = ""
)

func main() {
	app := kingpin.New("watch", "Like the unix one but works cross-platform without magic.").
		Interspersed(false)
	app.ErrorWriter(os.Stderr)

	w, err := watch.NewWatch(version, revision, built)
	if err != nil {
		app.Fatalf("%v\n", err)
	}

	w.ConfigureCli(app)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		app.Usage(nil)
		os.Exit(1)
	}
}
