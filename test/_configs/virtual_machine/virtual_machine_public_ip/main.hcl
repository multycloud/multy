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
  route_table_id = rt
  subnet_id      = subnet
}

multy "virtual_machine" "vm" {
  name      = "test-vm"
  os        = "linux"
  size      = "micro"
  user_data = cloud_specific_value({
    aws : "#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on AWS</h1>\" > /var/www/html/index.html",
    azure : "#!/bin/bash -xe\nsudo su\n yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on Azure</h1>\" > /var/www/html/index.html;",
  })
  subnet_id = subnet
  ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCwSjgjEIKewBWACaOVGg4qsGSHhIeteCmbtn4/DpL0yugLf5c/K/RJVQOKG+dVXVfWD3oAb4JY8jvkZdVACcuocoCewrIEHXxZJmGehxgCeUG8HZ+14mODosUOUYCe3kKCWU2SnUhjX+8x6btxqDOEhtghN3qR52kjm/OUw0Ap43weR1sdkJwtUz7CAXzdCxEKj16R0SY/dNn3uIISPetqm7vqy0ecMJdasbj/X6IAKeiZHe5UKtmOCGYMLwYfKqsrEnzk5rCfa3PK0iYoPv8AB3ocpONcuBshJyDZWhaFBhBrs5SGrWcF34wckD37SNtRZJt+Fuaxe8MpqUVueGTgFViKokCxCfbTnKWRbdGXpSfS6Q0OSvZTWkUEy5ZjxsA03LT4Bcbzq19sABdbyrcEMdv8bq0fhNyGJcGYNJr2uC4+J7irXAM/TuFje4CpJ0G+J3gCrQ2BUeWOYBdjfeP+LckgVXP+TMcEEe4iq5B9psyIS7o58KeNQdFH9jQteIE= joao@Joaos-MBP"
  public_ip = true
}