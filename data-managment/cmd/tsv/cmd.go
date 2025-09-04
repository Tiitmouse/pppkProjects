package tsv

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"go.uber.org/zap"
)

type TsvCmd struct{}

// Execute implements subcommands.Command.
func (t *TsvCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		zap.S().Errorf("No arguments for path to file")
		return subcommands.ExitUsageError
	}

	pathToFile := f.Args()[0]
	zap.S().Debugf("Path to file %s", pathToFile)

	err := action(pathToFile)
	if err != nil {
		return 1
	}

	return 0
}

// Name implements subcommands.Command.
func (t *TsvCmd) Name() string {
	return "parse"
}

// SetFlags implements subcommands.Command.
func (t *TsvCmd) SetFlags(*flag.FlagSet) {
	// TODO: seet flags if any
}

// Synopsis implements subcommands.Command.
func (t *TsvCmd) Synopsis() string {
	return "Parses the given tsv file"
}

// Usage implements subcommands.Command.
func (t *TsvCmd) Usage() string {
	return "parse ./path/to/file.tsv"
}
