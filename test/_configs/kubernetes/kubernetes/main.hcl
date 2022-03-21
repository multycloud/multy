config {
  location = "ireland"
  clouds   = ["aws", "azure"]
}

multy "kubernetes_service" "example" {
  name       = "example"
  subnet_ids = [subnet1, subnet2]
}

multy "kubernetes_node_pool" "example_pool" {
  name            = "example"
  cluster_id      = example
  subnet_ids      = [subnet1, subnet2]
  max_node_count  = 1
  min_node_count  = 1
  is_default_pool = true
  vm_size         = "medium"
}


multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet1" {
  name              = "private-subnet"
  cidr_block        = "10.0.1.0/24"
  virtual_network   = example_vn
  availability_zone = 1
}
multy "subnet" "subnet2" {
  name              = "public-subnet"
  cidr_block        = "10.0.2.0/24"
  virtual_network   = example_vn
  availability_zone = 2
}

multy route_table "rt" {
  name            = "test-rt"
  virtual_network = example_vn
  routes          = [
    {
      cidr_block  = "0.0.0.0/0"
      destination = "internet"
    }
  ]
}
multy route_table_association rta {
  route_table_id = rt
  subnet_id      = subnet2
}

output kubernetes_outputs {
  value = {
    aws_endpoint         = aws.example.endpoint,
    azure_endpoint       = azure.example.endpoint,
    aws_ca_certificate   = aws.example.ca_certificate,
    azure_ca_certificate = azure.example.ca_certificate
  }
}