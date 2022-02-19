package terraform

import (
	"fmt"
	"multy-go/resources/output"
)

const RandomPasswordResourceName = "random_password"

type RandomPassword struct {
	*output.TerraformResource `hcl:",squash" default:"name=random_password"`
	Length                    int  `hcl:"length"`
	Special                   bool `hcl:"special"`
	Upper                     bool `hcl:"upper"`
	Lower                     bool `hcl:"lower"`
	Number                    bool `hcl:"number"`
}

func (r *RandomPassword) GetResult() string {
	return fmt.Sprintf("random_password.%s.result", r.ResourceId)
}
