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
  virtual_network_id = example_vn.id
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