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

multy route_table "rt" {
  name               = "test-rt"
  virtual_network_id = example_vn.id
  routes             = [
    {
      cidr_block  = "0.0.0.0/0"
      destination = "internet"
    }
  ]
}

multy route_table_association rta {
  route_table_id = rt.id
  subnet_id      = subnet.id
}

multy "virtual_machine" "vm" {
  name      = "test-vm"
  os        = "linux"
  size      = "micro"
  user_data = "sudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html;"
  subnet_id = subnet.id
  ssh_key_file_name = "./ssh_key.pub"
  public_ip = true
}