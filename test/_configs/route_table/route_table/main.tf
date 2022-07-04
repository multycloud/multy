resource "aws_vpc" "example_vn_aws" {
  tags                 = { "Name" = "example_vn" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  provider             = "aws.eu-west-1"
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags     = { "Name" = "example_vn" }
  vpc_id   = aws_vpc.example_vn_aws.id
  provider = "aws.eu-west-1"
}
resource "aws_default_security_group" "example_vn_aws" {
  tags   = { "Name" = "example_vn" }
  vpc_id = aws_vpc.example_vn_aws.id
  ingress {
    protocol  = "-1"
    from_port = 0
    to_port   = 0
    self      = true
  }
  egress {
    protocol  = "-1"
    from_port = 0
    to_port   = 0
    self      = true
  }
  provider = "aws.eu-west-1"
}
resource "azurerm_virtual_network" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example_vn"
  location            = "northeurope"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example_vn"
  location            = "northeurope"
  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "google_compute_network" "example_vn_gcp" {
  name                            = "example_vn"
  project                         = "multy-project"
  routing_mode                    = "REGIONAL"
  description                     = "Managed by Multy"
  auto_create_subnetworks         = false
  delete_default_routes_on_create = true
  provider                        = "google.europe-west1"
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
}
resource "aws_route_table" "rt_aws" {
  tags   = { "Name" = "test-rt" }
  vpc_id = aws_vpc.example_vn_aws.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
  provider = "aws.eu-west-1"
}
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-rt"
  location            = "northeurope"
  route {
    name           = "internet"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "Internet"
  }
}
resource "google_compute_route" "rt_gcp-0" {
  name             = "test-rt-0"
  project          = "multy-project"
  dest_range       = "0.0.0.0/0"
  network          = google_compute_network.example_vn_gcp.id
  priority         = 1000
  tags             = ["subnet-subnet1", "subnet-subnet2"]
  next_hop_gateway = "default-internet-gateway"
  provider         = "google.europe-west1"
}
resource "aws_subnet" "subnet1_aws-1" {
  tags              = { "Name" = "subnet1-1" }
  cidr_block        = "10.0.1.0/25"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet1_aws-2" {
  tags              = { "Name" = "subnet1-2" }
  cidr_block        = "10.0.1.128/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet1_aws-3" {
  tags              = { "Name" = "subnet1-3" }
  cidr_block        = "10.0.1.192/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1c"
  provider          = "aws.eu-west-1"
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "subnet1_azure" {
  subnet_id      = azurerm_subnet.subnet1_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
resource "google_compute_subnetwork" "subnet1_gcp" {
  name                     = "subnet1"
  project                  = "multy-project"
  ip_cidr_range            = "10.0.1.0/24"
  network                  = google_compute_network.example_vn_gcp.id
  private_ip_google_access = true
  provider                 = "google.europe-west1"
}
resource "aws_subnet" "subnet2_aws-1" {
  tags              = { "Name" = "subnet2-1" }
  cidr_block        = "10.0.2.0/25"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet2_aws-2" {
  tags              = { "Name" = "subnet2-2" }
  cidr_block        = "10.0.2.128/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet2_aws-3" {
  tags              = { "Name" = "subnet2-3" }
  cidr_block        = "10.0.2.192/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1c"
  provider          = "aws.eu-west-1"
}
resource "azurerm_subnet" "subnet2_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet2"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "subnet2_azure" {
  subnet_id      = azurerm_subnet.subnet2_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
resource "google_compute_subnetwork" "subnet2_gcp" {
  name                     = "subnet2"
  project                  = "multy-project"
  ip_cidr_range            = "10.0.2.0/24"
  network                  = google_compute_network.example_vn_gcp.id
  private_ip_google_access = true
  provider                 = "google.europe-west1"
}
provider "aws" {
  region = "eu-west-1"
  alias  = "eu-west-1"
}
provider "azurerm" {
  features {
  }
}
provider "google" {
  region = "europe-west1"
  alias  = "europe-west1"
}
