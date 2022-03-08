package cli

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"multy-go/decoder"
	"multy-go/encoder"
	"multy-go/parser"
	"multy-go/resources"
	"path/filepath"
)

type CheckCommand struct {
	ConfigFiles []string
}

func (c *CheckCommand) ParseFlags(f *flag.FlagSet, args []string) {
	_ = f.Parse(args)
	c.ConfigFiles = f.Args()
}

func (c *CheckCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "check",
		Description: "validates the given config files, printing any errors found",
		Usage:       "multy check [files...] [options]",
	}
}

func (c *CheckCommand) Execute(ctx context.Context) error {
	if len(c.ConfigFiles) == 0 {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return fmt.Errorf("error while reading current directory: %s", err.Error())
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".mt" {
				log.Println("config file found:", file.Name())
				c.ConfigFiles = append(c.ConfigFiles, file.Name())
			}
		}
	}

	if len(c.ConfigFiles) == 0 {
		return fmt.Errorf("no .mt config files found")
	}

	p := parser.Parser{CliVars: nil}
	parsedConfig := p.Parse(c.ConfigFiles)

	r := decoder.Decode(parsedConfig)
	mctx := resources.MultyContext{Resources: r.Resources, Location: r.GlobalConfig.Location}

	_ = encoder.TranslateResources(r, mctx)

	fmt.Println("no validation errors found")
	return nil
}
