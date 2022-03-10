config {
  location = "us-east"
  clouds   = ["aws", "azure"]
}
multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet1" {
  name              = "subnet1"
  cidr_block        = "10.0.1.0/24"
  virtual_network   = example_vn
  availability_zone = 1
}
multy "subnet" "subnet2" {
  name              = "subnet2"
  cidr_block        = "10.0.2.0/24"
  virtual_network   = example_vn
  availability_zone = 2
}
multy "subnet" "subnet3" {
  name            = "subnet3"
  cidr_block      = "10.0.3.0/24"
  virtual_network = example_vn
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
  subnet_id      = subnet3.id
}
multy route_table_association rta2 {
  route_table_id = rt.id
  subnet_id      = subnet2.id
}
multy route_table_association rta3 {
  route_table_id = rt.id
  subnet_id      = subnet1.id
}
multy "database" "example_db" {
  name           = "example-db"
  size           = "nano"
  engine         = "mysql"
  engine_version = "5.7"
  storage        = 10
  db_username    = "multyadmin"
  db_password    = "multy-Admin123!"
  subnet_ids     = [
    subnet1.id,
    subnet2.id,
  ]
  clouds         = ["azure"]
}

multy "virtual_machine" "vm" {
  name              = "test-vm"
  os                = "linux"
  size              = "micro"
  user_data         = templatefile("init.sh", {
    db_host : azure.example_db.host,
    db_username : azure.example_db.username,
    db_password : azure.example_db.password
  })
  subnet_id         = subnet3.id
  ssh_key_file_path = "./ssh_key.pub"
  public_ip         = true

  depends_on = [example_db]
}
multy "network_security_group" nsg2 {
  name            = "test-nsg2"
  virtual_network = example_vn
  rules           = [
    {
      protocol   = "tcp"
      priority   = "100"
      action     = "allow"
      from_port  = "80"
      to_port    = "80"
      cidr_block = "0.0.0.0/0"
      direction  = "both"
    }, {
      protocol   = "tcp"
      priority   = "120"
      action     = "allow"
      from_port  = "22"
      to_port    = "22"
      cidr_block = "0.0.0.0/0"
      direction  = "both"
    }, {
      protocol   = "tcp"
      priority   = "140"
      action     = "allow"
      from_port  = "443"
      to_port    = "443"
      cidr_block = "0.0.0.0/0"
      direction  = "both"
    }
  ]
}
output "aws_endpoint" {
  value = "http://${aws.vm.public_ip}:4000"
}
output "azure_endpoint" {
  value = "http://${azure.vm.public_ip}:4000"
}
