//go:build e2e
// +build e2e

package database

import (
	"context"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"multy-go/cli"
	"os"
	"os/exec"
	"testing"
)

const DestroyAfter = false

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

var config = fmt.Sprintf(`
multy virtual_network vn {
  name       = "db-vn"
  cidr_block = "10.0.0.0/16"
}
multy subnet subnet1 {
  name              = "db-subnet1"
  cidr_block        = "10.0.0.0/24"
  virtual_network   = vn
  availability_zone = 1
}
multy subnet subnet2 {
  name              = "db-subnet2"
  cidr_block        = "10.0.1.0/24"
  virtual_network   = vn
  availability_zone = 2
}
multy route_table "rt" {
  name            = "db-rt"
  virtual_network = vn
  routes          = [
    {
      cidr_block  = "0.0.0.0/0"
      destination = "internet"
    }
  ]
}

multy route_table_association rta {
  route_table_id = rt.id
  subnet_id      = subnet1.id
}
multy route_table_association rta2 {
  route_table_id = rt.id
  subnet_id      = subnet2.id
}
multy "database" "db" {
  name           = "dbhlmzapdo"
  size           = "nano"
  engine         = "mysql"
  engine_version = "5.7"
  storage        = 10
  db_username    = "%s"
  db_password    = "%s"
  subnet_ids     = [
    subnet1.id,
    subnet2.id,
  ]
}
`, db_username, db_passwd)

func testDb(t *testing.T, cloudSpecificConfig string, cloudName string) {
	t.Parallel()

	multyFileName := fmt.Sprintf("database_%s.hcl", cloudName)
	tfDir := fmt.Sprintf("terraform_%s", cloudName)

	err := os.WriteFile(multyFileName, []byte(cloudSpecificConfig+config), 0664)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer os.Remove(multyFileName)

	cmd := cli.TranslateCommand{}
	cmd.OutputFile = tfDir + "/main.tf"

	ctx := context.Background()
	err = cmd.Execute([]string{multyFileName}, ctx)
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

	host := terraform.Output(t, tfOptions, "db_host")

	var username = db_username
	// TODO: fix this. azure and aws should have similar databases and able to be accessed the same way
	if cloudName == "azure" {
		username = fmt.Sprintf("%s@%s", db_username, host)
	}

	// TODO: fix Client with IP address '188.154.29.237' is not allowed to connect to this MySQL server in azure
	out, err := exec.Command("/usr/bin/mysql", "-h", host, "-P", "3306", "-u", username, "--password="+db_passwd, "-e", "select 12+34;").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	assert.Contains(t, string(out), "46")
}

func TestAwsDb(t *testing.T) {
	testDb(t, aws_global_config, "aws")
}

func TestAzureDb(t *testing.T) {
	testDb(t, azure_global_config, "azure")
}
