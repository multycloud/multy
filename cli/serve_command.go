package cli

import (
	"context"
	flag "github.com/spf13/pflag"
	"multy-go/api"
)

// ServeCommand - Temporary command to start the grpc server. Eventually we won't have a CLI / it will live in a diff repo.
type ServeCommand struct {
}

func (c *ServeCommand) ParseFlags(f *flag.FlagSet, args []string) {
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
	api.RunServer(ctx)
	return nil
}
