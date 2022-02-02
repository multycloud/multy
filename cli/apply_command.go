package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"io/ioutil"
	"os"
)

type ApplyCommand struct {
	Translate *TranslateCommand
}

func (c *ApplyCommand) Name() string {
	return "apply"
}

func (c *ApplyCommand) Init() {
	c.Translate = &TranslateCommand{}
	c.Translate.Init()
}

func (c *ApplyCommand) Execute(args []string, ctx context.Context) error {
	err := c.Translate.Execute(args, nil)

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
