package scrape

import (
	"context"
	"data-managment/util/bucket"
	"flag"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"go.uber.org/zap"
)

type ScrapeCmd struct{}

// Execute implements subcommands.Command.
func (t *ScrapeCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	dirPath := "./cmd/scrape/scrapper/dwnData"

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		zap.S().Errorf("Failed to read directory: %s", err)
		return 1
	}

	files := make([]*os.File, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dirPath, entry.Name())
			file, err := os.OpenFile(filePath, 0, 0755)
			if err != nil {
				zap.S().Errorf("Failed to read file %s, err = %v", filePath, err)
				continue
			}
			files = append(files, file)
			defer file.Close()
		}
	}
	count, err := bucket.Bucket.UploadMany(files)
	if err != nil {
		zap.S().Errorf("Failed with file uploading, err = %v", err)
		return 1
	}

	zap.S().Infof("Uploading Finished, %d files uploaded", count)
	return 0
}

// Name implements subcommands.Command.
func (t *ScrapeCmd) Name() string {
	return "scrape"
}

// SetFlags implements subcommands.Command.
func (t *ScrapeCmd) SetFlags(*flag.FlagSet) {
	// TODO: add flags if any
}

// Synopsis implements subcommands.Command.
func (t *ScrapeCmd) Synopsis() string {
	return "Uses a scrapper to scrape data from web"
}

// Usage implements subcommands.Command.
func (t *ScrapeCmd) Usage() string {
	return "scrape"
}
