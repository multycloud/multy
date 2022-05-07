//go:build e2e
// +build e2e

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/multycloud/multy/cmd"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"os"
	"os/exec"
	"testing"
)

const DestroyAfter = true

const aws_global_config = `
config {
  clouds   = ["aws"]
  location = "EU_WEST_1"
}
`

const azure_global_config = `
config {
  clouds   = ["azure"]
  location = "us-east"
}
`

var config = `
multy "kubernetes_service" "kubernetes_test" {
    name = "kbn-test"
    subnet_ids = [subnet1.id, subnet2.id]
}

multy "kubernetes_node_pool" "kbn_test_pool" {
  name = "kbntestpool"
  cluster_id = kubernetes_test.id
  subnet_ids = [subnet1.id, subnet2.id]
  starting_node_count = 1
  max_node_count = 1
  min_node_count = 1
  labels = { "multy.dev/env": "test" }
  vm_size = "medium"
  is_default_pool = true
}


multy "virtual_network" "kbn_test_vn" {
  name       = "kbn_test_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet1" {
  name              = "private-subnet"
  cidr_block        = "10.0.1.0/24"
  virtual_network   = kbn_test_vn
  availability_zone = 1
}
multy "subnet" "subnet2" {
  name              = "public-subnet"
  cidr_block        = "10.0.2.0/24"
  virtual_network   = kbn_test_vn
  availability_zone = 2
}

multy route_table "rt" {
  name            = "test-rt"
  virtual_network = kbn_test_vn
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
	Status string `json:"status,omitempty"`
}

type Node struct {
	Status struct {
		Conditions []NodeStatusCondition `json:"conditions,omitempty"`
	} `json:"status"`
	Metadata struct {
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

func testKubernetes(t *testing.T, cloudSpecificConfig string, cloudName string) {
	multyFileName := fmt.Sprintf("kubernetes_%s.hcl", cloudName)
	tfDir := fmt.Sprintf("terraform_%s", cloudName)

	err := os.WriteFile(multyFileName, []byte(cloudSpecificConfig+config), 0664)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer os.Remove(multyFileName)

	cmd := cmd.TranslateCommand{}
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	cmd.ParseFlags(f, []string{multyFileName})
	err = f.Set("output", tfDir+"/main.tf")
	if err != nil {
		t.Fatal(err.Error())
	}

	ctx := context.Background()
	err = cmd.Execute(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer func() {
		if DestroyAfter {
			os.RemoveAll(tfDir)
		}
	}()

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{TerraformDir: tfDir})
	terraform.InitAndApply(t, tfOptions)

	defer func() {
		if DestroyAfter {
			terraform.Destroy(t, tfOptions)
		}
	}()

	// update kubectl configuration so that we can use kubectl commands - probably can't run this in parallel
	if cloudName == "aws" {
		// aws eks --region eu-west-1 update-kubeconfig --name kubernetes_test
		out, err := exec.Command("/usr/bin/aws", "eks", "--region", "eu-west-1", "update-kubeconfig", "--name", "kbn-test").CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
		}
	} else {
		// az aks get-credentials --resource-group ks-rg --name example
		out, err := exec.Command("/usr/bin/az", "aks", "get-credentials", "--resource-group", "ks-rg", "--name", "kbn-test").CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
		}
	}
	// kubectl get nodes -o json
	out, err := exec.Command("/usr/local/bin/kubectl", "get", "nodes", "-o", "json").CombinedOutput()
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
		Status: "True",
	})

	labels := o.Items[0].Metadata.Labels
	assert.Contains(t, maps.Keys(labels), "multy.dev/env")
	assert.Equal(t, labels["multy.dev/env"], "test")
}

func TestAwsKubernetes(t *testing.T) {
	testKubernetes(t, aws_global_config, "aws")
}
func TestAzureKubernetes(t *testing.T) {
	testKubernetes(t, azure_global_config, "azure")
}
