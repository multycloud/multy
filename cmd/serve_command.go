package cmd

import (
	"context"
	"github.com/multycloud/multy/api"
	"github.com/multycloud/multy/flags"
	flag "github.com/spf13/pflag"
)

// ServeCommand - Temporary command to start the grpc server. Eventually we won't have a CLI / it will live in a diff repo.
type ServeCommand struct {
	Port int
}

func (c *ServeCommand) ParseFlags(f *flag.FlagSet, args []string) {
	f.IntVar(&c.Port, "port", 8000, "Port to run server on")
	f.BoolVar(&flags.DryRun, "dry_run", false, "If true, nothing will be deployed")
	_ = f.Parse(args)
}

func (c *ServeCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "serve",
		Description: "starts multy server",
		Usage:       "multy serve [options]",
	}
}

func (c *ServeCommand) Execute(ctx context.Context) error {
	api.RunServer(ctx, c.Port)
	return nil
}
