multy "virtual_network" "example_vn" {
  rg_vars    = {
    app = "backend"
  }
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
  clouds     = ["aws"]
}
multy "subnet" "example_subnet" {
  rg_vars            = {
    app = "backend"
  }
  name               = "example_subnet"
  cidr_block         = "10.0.0.0/24"
  clouds             = ["aws"]
  virtual_network = aws.example_vn
}
config {
  default_resource_group_name = "${resource_type}-rg"
  location                    = var.location
  clouds                      = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "ireland"
}