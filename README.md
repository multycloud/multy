<a href="https://s3.eu-west-1.amazonaws.com/www.multy.dev/logo2.pngutm_source=github.com">
    <img src="https://s3.eu-west-1.amazonaws.com/www.multy.dev/logo2.png" width="80">
</a>

<br/>
<br/>

**Multy is the easiest way to deploy multi cloud infrastructure**

Write cloud-agnostic config deployed across multiple clouds.

Let's try to deploy a simple virtual machine into AWS and Azure.

```hcl
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
}
multy "virtual_machine" "vm" {
  name      = "test-vm"
  os        = "linux"
  size      = "micro"
  subnet_id = subnet.id
}
```

If we were to deploy this using Terraform, we would first need to understand how services such as `aws_vpc`
and `azurerm_virtual_network` behave and how they differ. Then we would need to define the same infrastructure
configuration twice, one for AWS and another for Azure.


<details><summary>This is the equivalent terraform configuration</summary>
<p>

```hcl
// multy:     19 lines
// terraform: 132 lines
resource "aws_vpc" "example_vn_aws" {
  tags =  {
    Name = "example_vn"
  }

  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags =  {
    Name = "example_vn"
  }

  vpc_id = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  tags =  {
    Name = "example_vn"
  }

  vpc_id = aws_vpc.example_vn_aws.id

  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    self        = true
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    self        = true
  }
}
resource "aws_subnet" "subnet_aws" {
  tags =  {
    Name = "subnet"
  }

  cidr_block = "10.0.2.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
resource "aws_instance" "vm_aws" {
  tags =  {
    Name = "test-vm"
  }

  ami           = "ami-09d4a659cdd8677be"
  instance_type = "t2.nano"
  subnet_id     = aws_subnet.subnet_aws.id
}
resource "azurerm_virtual_network" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "northeurope"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "northeurope"

  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "azurerm_subnet" "subnet_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "subnet_azure" {
  subnet_id      = azurerm_subnet.subnet_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "northeurope"

  ip_configuration {
    name                          = "internal"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    primary                       = true
  }
}
resource "azurerm_linux_virtual_machine" "vm_azure" {
  resource_group_name   = azurerm_resource_group.vm-rg.name
  name                  = "test-vm"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm_azure.id]

  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }

  admin_username = "multyadmin"
  admin_password = "Multyadmin090#"

  source_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7_9-gen2"
    version   = "latest"
  }

  disable_password_authentication = false
}
resource "azurerm_resource_group" "vm-rg" {
  name     = "vm-rg"
  location = "northeurope"
}
resource "azurerm_resource_group" "vn-rg" {
  name     = "vn-rg"
  location = "northeurope"
}
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {}
}
```

</p>
</details>

With Multy, **you write once, and deploy anywhere**.

### Getting started

---

* **[Docs](https://multy.dev/docs/intro/)**: check our documentation page to get started
* **[Examples](https://multy.dev/docs/examples/)**: look over some common examples
* **[Resources](https://multy.dev/docs/resources/virtual_network/)**: currently supported resources

