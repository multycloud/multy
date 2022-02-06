package terraform

import (
	"fmt"
)

const RandomPasswordResourceName = "random_password"

type RandomPassword struct {
	ResourceName string `hcl:",key"`
	ResourceId   string `hcl:",key"`
	Length       int    `hcl:"length"`
	Special      bool   `hcl:"special"`
	Upper        bool   `hcl:"upper"`
	Lower        bool   `hcl:"lower"`
	Number       bool   `hcl:"number"`
}

func (r *RandomPassword) GetResult() string {
	return fmt.Sprintf("random_password.%s.result", r.ResourceId)
}
