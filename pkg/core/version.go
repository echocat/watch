package watch

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

func NewVersion(name, revision, built string) (*Version, error) {
	if built == "" {
		built = time.Now().Format(time.RFC3339)
	}
	builtTs, err := time.Parse(time.RFC3339, built)
	if err != nil {
		return nil, fmt.Errorf("illegal built value '%s': %w", built, err)
	}
	if revision == "" {
		if revision, err = randomRevision(builtTs); err != nil {
			return nil, err
		}
	}

	return &Version{
		Name:     name,
		Revision: revision,
		Built:    builtTs,
	}, nil
}

type Version struct {
	Name     string
	Revision string
	Built    time.Time
}

func (this *Version) ConfigureCli(cli *kingpin.Application) {
	cli.Flag("version", "Will print information about this version.").
		PreAction(func(action *kingpin.ParseContext) error {
			fmt.Printf(`watch
 Version:    %s
 Revision:   %s
 Built:      %v
 Go version: %s
 OS/Arch:    %s/%s
`,
				this.Name, this.Revision, this.Built, runtime.Version(), runtime.GOOS, runtime.GOARCH)

			os.Exit(0)
			return nil
		}).
		Bool()
}

func randomRevision(baseOn time.Time) (string, error) {
	b := make([]byte, sha1.Size)
	rng := rand.New(rand.NewSource(baseOn.UnixNano()))
	if n, err := rng.Read(b); err != nil {
		panic(err)
	} else if n < len(b) {
		return "", io.ErrShortBuffer
	}
	return hex.EncodeToString(b), nil
}
