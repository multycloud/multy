package cmd

import (
	"context"
	flag "github.com/spf13/pflag"
)

type CommandDesc struct {
	Name        string
	Description string
	Usage       string
}

type Command interface {
	Description() CommandDesc
	Execute(ctx context.Context) error
	ParseFlags(set *flag.FlagSet, strings []string)
}
