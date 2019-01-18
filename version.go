package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/alecthomas/kingpin"
	"io"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	_ = app.Flag("version", "Will print information about this version.").
		PreAction(func(action *kingpin.ParseContext) error {
			fmt.Printf(`watch
 Version:    %s
 Revision:   %s
 Built:      %v
 Go version: %s
 OS/Arch:    %s/%s
`,
				version, revision, builtTs, runtime.Version(), runtime.GOOS, runtime.GOARCH)

			os.Exit(0)
			return nil
		}).
		Bool()

	version  = "development"
	revision = ""
	built    = ""
	builtTs  time.Time
)

func init() {
	if built == "" {
		built = time.Now().Format(time.RFC3339)
	}
	var err error
	if builtTs, err = time.Parse(time.RFC3339, built); err != nil {
		panic(fmt.Sprintf("illegal built value '%s': %v", built, err))
	}
	if revision == "" {
		revision = randomRevision(builtTs)
	}
}

func randomRevision(baseOn time.Time) string {
	b := make([]byte, sha1.Size)
	rng := rand.New(rand.NewSource(baseOn.UnixNano()))
	if n, err := rng.Read(b); err != nil {
		panic(err)
	} else if n < len(b) {
		panic(io.ErrShortBuffer)
	}
	return hex.EncodeToString(b)
}
