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
	Env  string
}

func (c *ServeCommand) ParseFlags(f *flag.FlagSet, args []string) {
	f.IntVar(&c.Port, "port", 8000, "Port to run server on")
	f.BoolVar(&flags.DryRun, "dry_run", false, "If true, nothing will be deployed")
	f.BoolVar(&flags.NoTelemetry, "no_telemetry", false, "If true, no logs will be stored")
	f.StringVar(&c.Env, "env", "prod", "Environment - one of prod, dev or local")
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
	flags.Environment = flags.Env(c.Env)
	api.RunServer(ctx, c.Port)
	return nil
}
