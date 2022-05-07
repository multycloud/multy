package cmd

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
)

var Version = "dev"

type VersionCommand struct {
}

func (c *VersionCommand) ParseFlags(f *flag.FlagSet, args []string) {
	_ = f.Parse(args)
}

func (c *VersionCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "version",
		Description: "prints version of executable",
	}
}

func (c *VersionCommand) Execute(ctx context.Context) error {
	fmt.Println(Version)
	return nil
}
