package cli

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"multy-go/decoder"
	"multy-go/encoder"
	"multy-go/parser"
	"multy-go/variables"
	"os"
	"path/filepath"
	"time"
)

type TranslateCommand struct {
	CommandLineVars variables.CommandLineVariables
	OutputFile      string
}

func (c *TranslateCommand) Name() string {
	return "translate"
}

func (c *TranslateCommand) Init() {
	flag.Var(&c.CommandLineVars, "var", "Variables to be passed to configuration")
	flag.StringVar(&c.OutputFile, "output", "main.tf", "Name of the output file")
}

func (c *TranslateCommand) Execute(args []string, ctx context.Context) error {
	start := time.Now()
	configFiles := args

	if len(configFiles) == 0 {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return fmt.Errorf("error while reading current directory: %s", err.Error())
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".mt" {
				log.Println("config file found:", file.Name())
				configFiles = append(configFiles, file.Name())
			}
		}
	}

	if len(configFiles) == 0 {
		return fmt.Errorf("no .mt config files found")
	}

	p := parser.Parser{CliVars: c.CommandLineVars}
	parsedConfig := p.Parse(configFiles)

	r := decoder.Decode(parsedConfig)

	hclOutput := encoder.Encode(r)

	d1 := []byte(hclOutput)
	err := os.WriteFile(c.OutputFile, d1, 0644)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err.Error())
	}

	log.Printf("multy finished translating: %s\n", time.Since(start).Round(time.Second))
	log.Println("output file:", c.OutputFile)
	return nil
}
