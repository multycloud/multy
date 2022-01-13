config {
  location = "ireland"
  clouds   = ["aws", "azure"]
}
multy "virtual_network" "example_vn" {
  name       = "example_vn"
  cidr_block = "10.0.0.0/16"
}
multy "subnet" "subnet1" {
  name               = "subnet1"
  cidr_block         = "10.0.1.0/24"
  virtual_network_id = example_vn.id
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
  subnet_id      = subnet1.id
}
multy "virtual_machine" "vm2" {
  name                       = "test-vm2"
  os                         = "linux"
  size                       = "micro"
  user_data = cloud_specific_value({
    aws : "#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on AWS</h1>\" > /var/www/html/index.html",
    azure : "#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on Azure</h1>\" > /var/www/html/index.html",
  })
  subnet_id                  = subnet1.id
  network_security_group_ids = [nsg2.id]
  public_ip                  = true
}
multy "virtual_machine" "vm" {
  name      = "test-vm"
  os        = "linux"
  size      = "micro"
  user_data = cloud_specific_value({
    aws : "#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on AWS</h1>\" > /var/www/html/index.html",
    azure : "#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on Azure</h1>\" > /var/www/html/index.html",
  })
  subnet_id = subnet1.id
  public_ip = true
}
multy "network_security_group" nsg2 {
  name               = "test-nsg2"
  virtual_network_id = example_vn.id
  rules              = [
    {
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
