config {
  clouds   = ["aws", "azure"]
  location = "uk"
}
multy "vault" "example" {
  name = "dev-test-secret-multy"
}
multy "vault_secret" "api_key" {
  name  = "api-key"
  vault = example
  value = "xxx"
}
multy "vault_access_policy" "kv_ap" {
  vault    = example
  identity = vm.identity
  access   = "owner"
}
multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet" {
  name            = "subnet"
  cidr_block      = "10.0.2.0/24"
  virtual_network = example_vn
}
multy "virtual_machine" "vm" {
  name      = "test-vm"
  os        = "linux"
  size      = "micro"
  subnet_id = subnet
  public_ip = true
}
