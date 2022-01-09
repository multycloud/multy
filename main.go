package main

import (
	"flag"
	"fmt"
	"log"
	"multy-go/decoder"
	"multy-go/encoder"
	"multy-go/parser"
	"multy-go/variables"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var commandLineVars variables.CommandLineVariables

	flag.Var(&commandLineVars, "var", "Variables to be passed to configuration")
	outputFile := flag.String("output", "main.tf", "Name of the output file")
	flag.Parse()
	configFiles := flag.Args()

	if len(configFiles) == 0 {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || filepath.Ext(path) != ".mt" {
				return nil
			}
			configFiles = append(configFiles, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	if len(configFiles) == 0 {
		log.Fatalf("no .mt config files found")
	}

	if len(configFiles) > 1 {
		log.Fatalf("multiple config files are not yet supported")
	}

	p := parser.Parser{CliVars: commandLineVars}
	parsedConfig := p.Parse(configFiles)

	r := decoder.Decode(parsedConfig)

	hclOutput := encoder.Encode(r)

	fmt.Println(hclOutput)
	d1 := []byte(hclOutput)
	_ = os.WriteFile(*outputFile, d1, 0644)
	_ = exec.Command("terraform", "fmt", "terraform/")
}
