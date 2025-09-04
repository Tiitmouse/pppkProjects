package main

import (
	"context"
	"flag"
	"os"

	"data-managment/app"
	"data-managment/cmd/scrape"
	"data-managment/cmd/show"
	"data-managment/cmd/tsv"
	"data-managment/cmd/version"
	"data-managment/util/repo"

	"github.com/google/subcommands"
)

func init() {
	app.Setup()

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&version.VersionCmd{}, "")
	subcommands.Register(&tsv.TsvCmd{}, "parse")
	subcommands.Register(&show.ShowCmd{}, "show")
	subcommands.Register(&scrape.ScrapeCmd{}, "scrape")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	defer repo.RepoCli.Disconnect()
	os.Exit(int(subcommands.Execute(ctx)))
}
