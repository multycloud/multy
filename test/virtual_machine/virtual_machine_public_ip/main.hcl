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
  virtual_network = example_vn
  availability_zone  = 2
}

multy route_table "rt" {
  name               = "test-rt"
  virtual_network    = example_vn
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
  user_data = cloud_specific_value({
    aws : "#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on AWS</h1>\" > /var/www/html/index.html",
    azure : "#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on Azure</h1>\" > /var/www/html/index.html",
  })
  subnet_id = subnet.id
  ssh_key_file_path = "./ssh_key.pub"
  public_ip = true
}