//go:build e2e
// +build e2e

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"multy-go/cli"
	"os"
	"os/exec"
	"testing"
)

const DestroyAfter = true

const db_username = "multy"
const db_passwd = "passwd1778!"

const aws_global_config = `
config {
  clouds   = ["aws"]
  location = "ireland"
}
output "db_host" {
  value = aws.db.host
}
`

const azure_global_config = `
config {
  clouds   = ["azure"]
  location = "us-east"
}
output "db_host" {
  value = azure.db.host
}
`

var config = `
multy "kubernetes_service" "kubernetes_test" {
    name = "example"
    subnet_ids = [subnet1.id, subnet2.id]
}

multy "kubernetes_node_pool" "example_pool" {
  name = "example"
  cluster_name = example.name
  subnet_ids = [subnet1.id, subnet2.id]
  starting_node_count = 1
  max_node_count = 1
  min_node_count = 1
}


multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet1" {
  name              = "private-subnet"
  cidr_block        = "10.0.1.0/24"
  virtual_network   = example_vn
  availability_zone = 1
}
multy "subnet" "subnet2" {
  name              = "public-subnet"
  cidr_block        = "10.0.2.0/24"
  virtual_network   = example_vn
  availability_zone = 2
}

multy route_table "rt" {
  name            = "test-rt"
  virtual_network = example_vn
  routes          = [
    {
      cidr_block  = "0.0.0.0/0"
      destination = "internet"
    }
  ]
}
multy route_table_association rta {
  route_table_id = rt.id
  subnet_id      = subnet2.id
}
`

type GetNodesOutput struct {
	Items []Node `json:"items"`
}

type NodeStatusCondition struct {
	Type   string `json:"type,omitempty"`
	Status bool   `json:"status,omitempty"`
}

type Node struct {
	Status struct {
		Conditions []NodeStatusCondition `json:"conditions,omitempty"`
	} `json:"status"`
}

func testDb(t *testing.T, cloudSpecificConfig string, cloudName string) {
	t.Parallel()

	multyFileName := fmt.Sprintf("kubernetes%s.hcl", cloudName)
	tfDir := fmt.Sprintf("terraform_%s", cloudName)

	err := os.WriteFile(multyFileName, []byte(cloudSpecificConfig+config), 0664)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer os.Remove(multyFileName)

	cmd := cli.TranslateCommand{}
	cmd.OutputFile = tfDir + "/main.tf"
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	cmd.ParseFlags(f, []string{multyFileName})

	ctx := context.Background()
	err = cmd.Execute(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{TerraformDir: tfDir})
	terraform.Init(t, tfOptions)
	terraform.Apply(t, tfOptions)

	defer func() {
		if DestroyAfter {
			terraform.Destroy(t, tfOptions)
			os.RemoveAll(tfDir)
		}
	}()

	// aws eks --region eu-west-1 update-kubeconfig --name kubernetes_test
	out, err := exec.Command("/usr/bin/aws", "eks", "--region", "eu-west-1", "update-kubeconfig", "--name", "kubernetes_test").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}
	// kubectl get nodes -o json
	out, err = exec.Command("/usr/local/bin/kubectl", "get", "nodes", "-o", "json").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	o := GetNodesOutput{}

	err = json.Unmarshal(out, &o)
	if err != nil {
		t.Fatal(fmt.Errorf("output cant be parsed: %s", err.Error()))
	}

	assert.Len(t, o.Items, 1)
	assert.Contains(t, o.Items[0].Status.Conditions, NodeStatusCondition{
		Type:   "Ready",
		Status: true,
	})
}

func TestAwsDb(t *testing.T) {
	testDb(t, aws_global_config, "aws")
}

func TestAzureDb(t *testing.T) {
	testDb(t, azure_global_config, "azure")
}
