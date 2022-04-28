config {
  location = "EU_WEST_1"
  clouds   = ["aws", "azure"]
}
multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet1" {
  name               = "subnet1"
  cidr_block         = "10.0.1.0/24"
  virtual_network = example_vn
}
multy "subnet" "subnet2" {
  name               = "subnet2"
  cidr_block         = "10.0.2.0/24"
  virtual_network = example_vn
  availability_zone  = 2
}