package cmd

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
)

type HelpCommand struct {
	AvailableCommands []Command
}

func (c *HelpCommand) ParseFlags(f *flag.FlagSet, args []string) {
	_ = f.Parse(args)
}

func (c *HelpCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "help",
		Description: "prints all available commands",
	}
}

func (c *HelpCommand) Execute(ctx context.Context) error {
	fmt.Println("Usage: multy <command> [args]")
	fmt.Println("Available commands:")
	for _, c := range c.AvailableCommands {
		fmt.Printf("%-12s\t%s\n", c.Description().Name, c.Description().Description)
	}
	return nil
}
