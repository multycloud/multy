config {
  clouds                      = ["aws", "azure"]
  location                    = "us-east"
  default_resource_group_name = "dbt-${resource_type}-rg"
}
multy virtual_network vn {
  name       = "db-vn"
  cidr_block = "10.0.0.0/16"
}
multy subnet subnet1 {
  name              = "subnet1"
  cidr_block        = "10.0.0.0/24"
  virtual_network   = vn
  availability_zone = 1
}
multy subnet subnet2 {
  name              = "subnet2"
  cidr_block        = "10.0.1.0/24"
  virtual_network   = vn
  availability_zone = 2
}
multy route_table "rt" {
  name            = "db-rt"
  virtual_network = vn
  routes          = [
    {
      cidr_block  = "0.0.0.0/0"
      destination = "internet"
    }
  ]
}

multy route_table_association rta {
  route_table_id = rt.id
  subnet_id      = subnet1.id
}
multy route_table_association rta2 {
  route_table_id = rt.id
  subnet_id      = subnet2.id
}
multy "database" "example_db" {
  name           = "example-db"
  size           = "nano"
  engine         = "mysql"
  engine_version = "5.7"
  storage        = 10
  db_username    = "multyadmin"
  db_password    = "multy$Admin123!"
  subnet_ids     = [
    subnet1.id,
    subnet2.id,
  ]
}
