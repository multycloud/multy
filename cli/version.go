package cli

import (
	"context"
	"fmt"
)

var Version = "dev"

type VersionCommand struct {
}

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Init() {
}

func (c *VersionCommand) Execute(args []string, ctx context.Context) error {
	fmt.Println(Version)
	return nil
}
