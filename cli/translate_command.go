package cli

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/proto/creds"
	"github.com/multycloud/multy/decoder"
	"github.com/multycloud/multy/encoder"
	"github.com/multycloud/multy/parser"
	"github.com/multycloud/multy/validate"
	"github.com/multycloud/multy/variables"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type TranslateCommand struct {
	CommandLineVars variables.CommandLineVariables
	OutputFile      string
	ConfigFiles     []string
}

func (c *TranslateCommand) ParseFlags(f *flag.FlagSet, args []string) {
	f.Var(&c.CommandLineVars, "var", "Variables to be passed to configuration")
	f.StringVar(&c.OutputFile, "output", "main.tf", "Name of the output file")
	_ = f.Parse(args)
	c.ConfigFiles = f.Args()
}

func (c *TranslateCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "translate",
		Description: "translates the multy configuration file(s) to a terraform file",
		Usage:       "multy translate [files...] [options]",
	}
}

func (c *TranslateCommand) Execute(ctx context.Context) error {
	start := time.Now()

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

	p := parser.Parser{CliVars: c.CommandLineVars}
	parsedConfig := p.Parse(c.ConfigFiles)

	r := decoder.Decode(parsedConfig)

	hclOutput, errs, err := encoder.Encode(r, &creds.CloudCredentials{})
	if len(errs) > 0 {
		validate.PrintAllAndExit(errs)
	}
	if err != nil {
		validate.LogInternalError(err.Error())
	}

	err = os.MkdirAll(filepath.Dir(c.OutputFile), os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err.Error())
	}
	err = os.WriteFile(c.OutputFile, []byte(hclOutput), os.ModePerm&0664)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err.Error())
	}

	fmt.Printf("multy finished translating: %s\n", time.Since(start).Round(time.Second))
	fmt.Println("output file:", c.OutputFile)
	return nil
}
