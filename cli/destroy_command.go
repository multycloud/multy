package cli

import (
	"context"
)

type DestroyCommand struct {
}

func (d *DestroyCommand) Name() string {
	return "destroy"
}

func (d *DestroyCommand) Init() {
}

func (d *DestroyCommand) Execute(args []string) error {
	tf, err := InitTf()
	if err != nil {
		return err
	}

	err = tf.Destroy(context.Background())
	if err != nil {
		return err
	}

	return nil
}
