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
  #  clouds         = ["azure"]
}

multy "virtual_machine" "vm" {
  name              = "test-vm"
  os                = "linux"
  size              = "micro"
  user_data         = cloud_specific_value({
    aws : templatefile("init.sh", {
      db_host : azure.example_db.host,
      db_username : azure.example_db.username,
      db_password : azure.example_db.password
    }),
    azure : templatefile("azure_init.sh", {
      vault_name : web_app_vault.name,
      db_host_secret_name : db_host.name,
      db_username_secret_name : db_username.name,
      db_password_secret_name : db_password.name
    })
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
multy "vault" "web_app_vault" {
  name = "web-app-vault-test"
}
multy "vault_secret" "db_host" {
  name  = "db-host"
  vault = web_app_vault
  value = azure.example_db.host
}
multy "vault_secret" "db_username" {
  name  = "db-username"
  vault = web_app_vault
  value = azure.example_db.username
}
multy "vault_secret" "db_password" {
  name  = "db-password"
  vault = web_app_vault
  value = azure.example_db.password
}
# https://stackoverflow.com/questions/68018819/reuse-iam-instance-profile
multy "vault_access_policy" "kv_ap" {
  vault    = web_app_vault
  identity = vm.identity
  access   = "owner"
}
output "aws_endpoint" {
  value = "http://${aws.vm.public_ip}:4000"
}
output "azure_endpoint" {
  value = "http://${azure.vm.public_ip}:4000"
}
