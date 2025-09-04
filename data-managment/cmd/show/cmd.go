package show

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"go.uber.org/zap"
)

type ShowCmd struct{}

// Execute implements subcommands.Command.
func (t *ShowCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if len(f.Args()) != 1 {
		zap.S().Errorf("No arguments for patient code")
		return subcommands.ExitUsageError
	}

	patientCode := f.Args()[0]
	zap.S().Debugf("patientCode %s", patientCode)

	err := connectData(patientCode)
	if err != nil {
		return 1
	}

	return 0
}

// Name implements subcommands.Command.
func (t *ShowCmd) Name() string {
	return "show"
}

// SetFlags implements subcommands.Command.
func (t *ShowCmd) SetFlags(*flag.FlagSet) {
	// TODO: add flags if any
}

// Synopsis implements subcommands.Command.
func (t *ShowCmd) Synopsis() string {
	return "Displays data for tsv file"
}

// Usage implements subcommands.Command.
func (t *ShowCmd) Usage() string {
	return "show PATIENT-CODE"
}
