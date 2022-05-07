package cmd

import (
	"context"
	"fmt"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
)

type printfer struct {
}

func (p printfer) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	fmt.Println()
}
func InitTf(ctx context.Context) (*tfexec.Terraform, error) {
	installer := &releases.LatestVersion{
		Product: product.Terraform,
	}

	execPath, err := installer.Install(ctx)
	if err != nil {
		return nil, fmt.Errorf("error installing Terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(".", execPath)
	if err != nil {
		return nil, fmt.Errorf("error running NewTerraform: %s", err)
	}

	// TODO: hide this behind a flag?
	tf.SetLogger(printfer{})

	tf.SetStdout(os.Stdout)
	return tf, nil
}
