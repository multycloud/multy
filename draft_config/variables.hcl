# you can combine network_interface.hcl and azure.hcl with {cloud} block
azure {
  variable "azure_auth" {
    type = object({
      client_id       = string
      client_secret   = string
      subscription_id = string
      tenant_id       = string
    })
    default = {}
  }

  variable "cidr_block" {
    type = string
    default = "10.1.0.0/16"
  }

  variable resource_group_schema {
    schema = "{env}-{domain}"
  }
}

aws {
  variable "auth" {
    type = object({
      access_key        = string
      secret_access_key = string
    })
    default = {}
  }

  variable "cidr_block" {
    type = string
    default = "10.0.0.0/16"
  }
}

