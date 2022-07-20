resource "google_sql_database_instance" "example_db_GCP" {
  name                = "example-db"
  project             = "multy-project"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier              = "db-f1-micro"
    availability_type = "ZONAL"
    disk_autoresize   = false
    disk_size         = 10
    ip_configuration {
      authorized_networks {
        value = "0.0.0.0/0"
      }
    }
  }
  provider = "google.us-east4"
}
resource "google_sql_user" "example_db_GCP" {
  name     = "multyadmin"
  instance = google_sql_database_instance.example_db_GCP.name
  password = "multy$Admin123!"
  provider = "google.us-east4"
  project  = "multy-project"
}

resource "google_compute_firewall" "vn_GCP" {
  name          = "db-vn-default-allow-ingress"
  project       = "multy-project"
  network       = google_compute_network.vn_GCP.id
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  priority      = 65534
  allow {
    protocol = "all"
  }
  target_tags = ["vn-db-vn-default-nsg"]
  provider    = "google.us-east4"
}
resource "aws_db_subnet_group" "example_db_aws" {
  tags        = { "Name" = "example-db" }
  name        = "example-db"
  description = "Managed by Multy"
  subnet_ids  = [aws_subnet.subnet_aws-1.id, aws_subnet.subnet_aws-2.id, aws_subnet.subnet_aws-3.id]
  provider    = "aws.us-east-1"
}
resource "aws_security_group" "example_db_aws" {
  tags        = { "Name" = "example-db" }
  vpc_id      = aws_vpc.vn_aws.id
  name        = "example-db"
  description = "Default security group of example-db"
  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
  provider = "aws.us-east-1"
}
resource "aws_db_instance" "example_db_aws" {
  tags                   = { "Name" = "exampledb" }
  allocated_storage      = 10
  engine                 = "mysql"
  engine_version         = "5.7"
  username               = "multyadmin"
  password               = "multy$Admin123!"
  instance_class         = "db.t2.micro"
  identifier             = "example-db"
  skip_final_snapshot    = true
  db_subnet_group_name   = aws_db_subnet_group.example_db_aws.name
  publicly_accessible    = true
  vpc_security_group_ids = [aws_security_group.example_db_aws.id]
  provider               = "aws.us-east-1"
}
resource "azurerm_mysql_server" "example_db_azure" {
  resource_group_name              = azurerm_resource_group.rg1.name
  name                             = "example-db"
  location                         = "eastus"
  administrator_login              = "multyadmin"
  administrator_login_password     = "multy$Admin123!"
  sku_name                         = "GP_Gen5_2"
  storage_mb                       = 10240
  version                          = "5.7"
  ssl_enforcement_enabled          = false
  ssl_minimal_tls_version_enforced = "TLSEnforcementDisabled"
}
resource "azurerm_mysql_virtual_network_rule" "example_db_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example-db"
  server_name         = azurerm_mysql_server.example_db_azure.name
  subnet_id           = azurerm_subnet.subnet_azure.id
}
resource "azurerm_mysql_firewall_rule" "example_db_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "public"
  server_name         = azurerm_mysql_server.example_db_azure.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "eastus"
}
resource "aws_route_table" "rt_aws" {
  tags   = { "Name" = "db-rt" }
  vpc_id = aws_vpc.vn_aws.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vn_aws.id
  }
  provider = "aws.us-east-1"
}
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-rt"
  location            = "eastus"
  route {
    name           = "internet"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "Internet"
  }
}
resource "aws_route_table_association" "rta_aws-1" {
  subnet_id      = aws_subnet.subnet_aws-1.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.us-east-1"
}
resource "aws_route_table_association" "rta_aws-2" {
  subnet_id      = aws_subnet.subnet_aws-2.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.us-east-1"
}
resource "aws_route_table_association" "rta_aws-3" {
  subnet_id      = aws_subnet.subnet_aws-3.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.us-east-1"
}
resource "azurerm_subnet_route_table_association" "subnet_azure" {
  subnet_id      = azurerm_subnet.subnet_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "google_compute_subnetwork" "subnet_GCP" {
  name                     = "subnet"
  project                  = "multy-project"
  ip_cidr_range            = "10.0.0.0/24"
  network                  = google_compute_network.vn_GCP.id
  private_ip_google_access = true
  provider                 = "google.us-east4"
}
resource "aws_subnet" "subnet_aws-1" {
  tags              = { "Name" = "subnet-1" }
  cidr_block        = "10.0.0.0/25"
  vpc_id            = aws_vpc.vn_aws.id
  availability_zone = "us-east-1a"
  provider          = "aws.us-east-1"
}
resource "aws_subnet" "subnet_aws-2" {
  tags              = { "Name" = "subnet-2" }
  cidr_block        = "10.0.0.128/26"
  vpc_id            = aws_vpc.vn_aws.id
  availability_zone = "us-east-1b"
  provider          = "aws.us-east-1"
}
resource "aws_subnet" "subnet_aws-3" {
  tags              = { "Name" = "subnet-3" }
  cidr_block        = "10.0.0.192/26"
  vpc_id            = aws_vpc.vn_aws.id
  availability_zone = "us-east-1c"
  provider          = "aws.us-east-1"
}
resource "azurerm_subnet" "subnet_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet"
  address_prefixes     = ["10.0.0.0/24"]
  virtual_network_name = azurerm_virtual_network.vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "google_compute_network" "vn_GCP" {
  name                            = "db-vn"
  project                         = "multy-project"
  routing_mode                    = "REGIONAL"
  description                     = "Managed by Multy"
  auto_create_subnetworks         = false
  delete_default_routes_on_create = true
  provider                        = "google.us-east4"
}
resource "aws_vpc" "vn_aws" {
  tags                 = { "Name" = "db-vn" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  provider             = "aws.us-east-1"
}
resource "aws_internet_gateway" "vn_aws" {
  tags     = { "Name" = "db-vn" }
  vpc_id   = aws_vpc.vn_aws.id
  provider = "aws.us-east-1"
}
resource "aws_default_security_group" "vn_aws" {
  tags   = { "Name" = "db-vn" }
  vpc_id = aws_vpc.vn_aws.id
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
  provider = "aws.us-east-1"
}
resource "azurerm_virtual_network" "vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-vn"
  location            = "eastus"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-vn"
  location            = "eastus"
  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
provider "aws" {
  region = "us-east-1"
  alias  = "us-east-1"
}
provider "azurerm" {
  features {
  }
}
provider "google" {
  region = "us-east4"
  alias  = "us-east4"
}
