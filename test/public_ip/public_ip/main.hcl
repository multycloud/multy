config {
  location = "ireland"
  clouds   = ["aws", "azure"]
}
multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet" {
  name               = "subnet"
  cidr_block         = "10.0.2.0/24"
  virtual_network_id = example_vn.id
  availability_zone  = 2
}
multy "network_interface" "public-nic" {
  name      = "test-public-nic"
  subnet_id = subnet.id
}
multy "network_interface" "private-nic" {
  name      = "test-private-nic"
  subnet_id = subnet.id
}
multy "public_ip" "ip" {
  name                 = "test-ip"
  network_interface_id = public-nic.id
}