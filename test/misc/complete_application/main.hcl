multy "virtual_network" "example_vn" {
  rg_vars    = {
    app = "backend"
  }
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "example_subnet" {
  rg_vars            = {
    app = "backend"
  }
  name               = "example_subnet"
  cidr_block         = "10.0.0.0/24"
  virtual_network = example_vn
}
multy "public_ip" "example_ip" {
  name = "example_ip"
}
multy "object_storage" "obj_storage" {
  name          = "test-storage"
  random_suffix = false
  location      = "uk"
}
multy "network_interface" "example_nic" {
  name      = "example_nic"
  subnet_id = example_subnet.id
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