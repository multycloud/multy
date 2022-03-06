config {
  location = "ireland"
  clouds   = ["aws"]
}

multy "kubernetes_service" "example" {
  name       = "example"
  subnet_ids = [subnet1.id, subnet2.id]
}

multy "kubernetes_node_pool" "example_pool" {
  name                = "example"
  cluster_name        = example.name
  subnet_ids          = [subnet1.id, subnet2.id]
  starting_node_count = 1
  max_node_count      = 1
  min_node_count      = 1
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
  route_table_id = rt.id
  subnet_id      = subnet2.id
}