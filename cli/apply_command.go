package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"os"
)

type ApplyCommand struct {
	Translate *TranslateCommand
}

func (c *ApplyCommand) ParseFlags(f *flag.FlagSet, args []string) {
	c.Translate = &TranslateCommand{}
	c.Translate.ParseFlags(f, args)
}

func (c *ApplyCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "apply",
		Description: "deploys or updates infrastructure according to the config file(s)",
		Usage:       "multy apply [files...] [options]",
	}
}

func (c *ApplyCommand) Execute(ctx context.Context) error {
	err := c.Translate.Execute(nil)

	if err != nil {
		return err
	}

	tf, err := InitTf(ctx)
	if err != nil {
		return err
	}

	tf.SetStdout(ioutil.Discard)

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		return fmt.Errorf("error running Init: %s", err)
	}

	tf.SetStdout(os.Stdout)

	err = tf.Apply(ctx)
	if err != nil {
		return fmt.Errorf("error running Apply: %s", err)
	}

	return nil

}
