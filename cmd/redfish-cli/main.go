package main

import (
	"github.com/alecthomas/kong"
)

type globalOptions struct {
	Verbose bool `help:"Enable verbose logging"`
}

type pklOptions struct {
	Pkl string `type:"path" help:"Path to PKL binary."`
}

var cli struct {
	globalOptions

	Update Update `cmd:"" help:"Update config for project."`
	Run    Run    `cmd:"" help:"Run locally using LocalRunManager."`
}

func main() {
	k := kong.Parse(&cli,
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			// Compact: true,
		}))

	err := k.Run(&cli.globalOptions)
	k.FatalIfErrorf(err)
}
