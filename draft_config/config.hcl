# Default auth env vars known to TF
# AWS
# export AWS_ACCESS_KEY_ID=""
# export AWS_SECRET_ACCESS_KEY=""
# --
# GCP
# export GOOGLE_APPLICATION_CREDENTIALS=""
# --
# Azure
# export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
# export ARM_CLIENT_SECRET="00000000-0000-0000-0000-000000000000"
# export ARM_SUBSCRIPTION_ID="00000000-0000-0000-0000-000000000000"
# export ARM_TENANT_ID="00000000-0000-0000-0000-000000000000"

variable "gcp_auth" {
  type    = object({
    application_credentials = string
  })
  default = {}
}

variable "azure_auth" {
  type    = object({
    client_id       = string
    client_secret   = string
    subscription_id = string
    tenant_id       = string
  })
  default = {}
}

config {
  clouds = ['aws', 'gcp']
}

resource virtual_network "example_network" {
  cidr_block = "10.0.0.0/16"
  name       = "main-{resource_type}" # network_interface: main-vpc; azure: main-vnet

  azure {
    resource_group = "{env}-main-vnet" # will override azure.resource_group_schema.schema
  }
}

resource subnet "example_subnet" {
  virtual_network = virtual_network.example_network.id
  name            = "main-{resource_type}" # for network_interface you cannot name a subnet, it will add it to tags
  cidr_block      = var.cidr_block # different values will be used based on the vendor
  cidr_blocks     = {
    aws : aws("vm.example_vm").ip
    az : "vm.example_vm.ip_address"
  }
}

aws {
  # resource to only be deployed in network_interface
  resource virtual_network_peering "vpc_peering" {
    from_virtual_network = virtual_network.vn1
    to_virtual_network   = virtual_network.vn2
  }

  terraform "aws_vpc" "main" {
    # terraform type to call terraform directly
    cidr_block       = var.cidr_block
    instance_tenancy = "default"

    tags = {
      Name = "main"
    }
  }
}

// change cloud provider
// deploy different resources to different clouds