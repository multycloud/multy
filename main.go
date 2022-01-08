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
)

func main() {
	log.SetFlags(log.Lshortfile)

	var commandLineVars variables.CommandLineVariables

	flag.Var(&commandLineVars, "var", "Variables to be passed to configuration")
	flag.Parse()
	configFiles := flag.Args()

	if len(configFiles) == 0 {
		log.Fatalf("no config file was specified")
	}

	if len(configFiles) > 1 {
		log.Fatalf("multiple config files are not yet supported")
	}

	log.Printf("parsed vars: %v", commandLineVars.String())

	p := parser.Parser{CliVars: commandLineVars}
	parsedConfig := p.Parse(configFiles[0])

	r := decoder.Decode(parsedConfig)

	hclOutput := encoder.Encode(r)

	//output aws_ip {
	//  value = aws_instance.example_vm_aws.public_ip
	//}
	//
	//output az_ip {
	//  value = azurerm_linux_virtual_machine.example_vm_azure.public_ip_address
	//}
	//
	//output az_storage_url {
	//  value = azurerm_storage_blob.file1_azure.url
	//}
	//
	//output aws_storage_url {
	//  value = "${aws_s3_bucket.storage_aws.bucket_domain_name}/${aws_s3_bucket_object.file1_aws.id}"
	//}

	fmt.Println(hclOutput)
	d1 := []byte(hclOutput)
	_ = os.WriteFile("terraform/main.tf", d1, 0644)
	_ = exec.Command("terraform", "fmt", "terraform/")
}
