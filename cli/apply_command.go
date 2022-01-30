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

func (c *ApplyCommand) Execute(args []string) error {
	err := c.Translate.Execute(args)

	if err != nil {
		return err
	}

	tf, err := InitTf()
	if err != nil {
		return err
	}

	tf.SetStdout(ioutil.Discard)

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return fmt.Errorf("error running Init: %s", err)
	}

	tf.SetStdout(os.Stdout)

	err = tf.Apply(context.Background())
	if err != nil {
		return fmt.Errorf("error running Apply: %s", err)
	}

	return nil

}
