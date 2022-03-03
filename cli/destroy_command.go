package cli

import (
	"context"
	flag "github.com/spf13/pflag"
)

type DestroyCommand struct {
}

func (d *DestroyCommand) ParseFlags(f *flag.FlagSet, args []string) {
	_ = f.Parse(args)
}

func (d *DestroyCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "destroy",
		Description: "destroys all infrastructure specified in the state file",
		Usage:       "multy destroy",
	}
}

func (d *DestroyCommand) Execute(ctx context.Context) error {
	tf, err := InitTf(ctx)
	if err != nil {
		return err
	}

	err = tf.Destroy(ctx)
	if err != nil {
		return err
	}

	return nil
}
