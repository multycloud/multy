package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
)

func InitTf() (*tfexec.Terraform, error) {
	installer := &releases.LatestVersion{
		Product: product.Terraform,
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error installing Terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(".", execPath)
	if err != nil {
		return nil, fmt.Errorf("error running NewTerraform: %s", err)
	}

	tf.SetStdout(os.Stdout)
	return tf, nil
}
