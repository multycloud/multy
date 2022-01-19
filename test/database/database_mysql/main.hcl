config {
  clouds = ["aws", "azure"]
  location = "uk"
}
multy virtual_network vn {
  name = "vn"
  cidr_block = "10.0.0.0/16"
}
multy subnet subnet1 {
  name = "subnet1"
  cidr_block = "10.0.0.0/24"
  virtual_network = vn
}
multy subnet subnet2 {
  name = "subnet2"
  cidr_block = "10.0.0.0/24"
  virtual_network = vn
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
