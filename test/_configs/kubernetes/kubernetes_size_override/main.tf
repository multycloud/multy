resource "aws_iam_role" "cluster_aws_default_pool" {
  tags               = { "Name" = "multy-k8nodepool-cluster_aws-node_pool_aws-role" }
  name               = "multy-k8nodepool-cluster_aws-node_pool_aws-role"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  provider           = "aws.eu-west-1"
}
resource "aws_iam_role_policy_attachment" "cluster_aws_default_pool_AmazonEKSWorkerNodePolicy" {
  role       = aws_iam_role.cluster_aws_default_pool.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  provider   = "aws.eu-west-1"
}
resource "aws_iam_role_policy_attachment" "cluster_aws_default_pool_AmazonEKS_CNI_Policy" {
  role       = aws_iam_role.cluster_aws_default_pool.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  provider   = "aws.eu-west-1"
}
resource "aws_iam_role_policy_attachment" "cluster_aws_default_pool_AmazonEC2ContainerRegistryReadOnly" {
  role       = aws_iam_role.cluster_aws_default_pool.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  provider   = "aws.eu-west-1"
}
resource "aws_eks_node_group" "cluster_aws_default_pool" {
  cluster_name    = aws_eks_cluster.cluster_aws.id
  node_group_name = "node_pool_aws"
  node_role_arn   = aws_iam_role.cluster_aws_default_pool.arn
  subnet_ids      = [
    aws_subnet.public_subnet_aws-1.id, aws_subnet.public_subnet_aws-2.id, aws_subnet.public_subnet_aws-3.id
  ]
  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }
  instance_types = ["t3.medium"]
  provider       = "aws.eu-west-1"
}
resource "aws_subnet" "cluster_aws_public_subnet" {
  tags              = { "Name" = "cluster_aws_public_subnet" }
  cidr_block        = "10.0.255.240/28"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "cluster_aws_private_subnet" {
  tags              = { "Name" = "cluster_aws_private_subnet" }
  cidr_block        = "10.0.255.224/28"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
  provider          = "aws.eu-west-1"
}
resource "aws_route_table" "cluster_aws_public_rt" {
  tags   = { "Name" = "cluster_aws_public_rt" }
  vpc_id = aws_vpc.example_vn_aws.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
  provider = "aws.eu-west-1"
}
resource "aws_route_table_association" "cluster_aws_public_rta" {
  subnet_id      = aws_subnet.cluster_aws_public_subnet.id
  route_table_id = aws_route_table.cluster_aws_public_rt.id
  provider       = "aws.eu-west-1"
}
resource "aws_iam_role" "cluster_aws" {
  tags               = { "Name" = "multy-k8cluster-cluster_aws-role" }
  name               = "multy-k8cluster-cluster_aws-role"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"eks.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  provider           = "aws.eu-west-1"
}
resource "aws_iam_role_policy_attachment" "cluster_aws_AmazonEKSClusterPolicy" {
  role       = aws_iam_role.cluster_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  provider   = "aws.eu-west-1"
}
resource "aws_iam_role_policy_attachment" "cluster_aws_AmazonEKSVPCResourceController" {
  role       = aws_iam_role.cluster_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  provider   = "aws.eu-west-1"
}
resource "aws_eks_cluster" "cluster_aws" {
  depends_on = [
    aws_subnet.cluster_aws_public_subnet, aws_subnet.cluster_aws_private_subnet, aws_route_table.cluster_aws_public_rt,
    aws_route_table_association.cluster_aws_public_rta,
    aws_iam_role_policy_attachment.cluster_aws_AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.cluster_aws_AmazonEKSVPCResourceController
  ]
  tags     = { "Name" = "cluster_aws" }
  role_arn = aws_iam_role.cluster_aws.arn
  vpc_config {
    subnet_ids              = [aws_subnet.cluster_aws_public_subnet.id, aws_subnet.cluster_aws_private_subnet.id]
    endpoint_private_access = true
  }
  kubernetes_network_config {
    service_ipv4_cidr = "10.100.0.0/16"
  }
  name     = "cluster_aws"
  provider = "aws.eu-west-1"
}
resource "azurerm_kubernetes_cluster" "cluster_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "cluster_azure"
  location            = "northeurope"
  default_node_pool {
    name                = "default"
    node_count          = 1
    max_count           = 1
    min_count           = 1
    enable_auto_scaling = true
    vm_size             = "Standard_DS2_v2"
    vnet_subnet_id      = azurerm_subnet.public_subnet_azure.id
  }
  dns_prefix = "clusterazureaks7hut"
  identity {
    type = "SystemAssigned"
  }
  network_profile {
    network_plugin     = "azure"
    dns_service_ip     = "10.100.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
    service_cidr       = "10.100.0.0/16"
  }
}
resource "google_service_account" "cluster_gcp" {
  account_id   = "cluster-gcp-clus5a4y-sa-mgby"
  display_name = "Service Account for cluster cluster-gcp - created by Multy"
  provider     = "google.europe-west1"
}
resource "google_container_node_pool" "cluster_gcp_default_pool" {
  name               = "node-pool-gcp"
  cluster            = google_container_cluster.cluster_gcp.id
  initial_node_count = 1
  node_locations     = ["europe-west1-d"]
  autoscaling {
    min_node_count = 1
    max_node_count = 1
  }
  node_config {
    machine_type    = "t2d-standard-2"
    tags            = ["subnet-public-subnet"]
    service_account = google_service_account.cluster_gcp.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]
  }
  provider = "google.europe-west1"
}
resource "google_container_cluster" "cluster_gcp" {
  name                     = "cluster-gcp"
  remove_default_node_pool = true
  initial_node_count       = 1
  subnetwork               = google_compute_subnetwork.public_subnet_gcp.id
  network                  = google_compute_network.example_vn_gcp.id
  ip_allocation_policy {
    services_ipv4_cidr_block = "10.100.0.0/16"
  }
  location = "europe-west1"
  node_config {
    machine_type    = "e2-micro"
    tags            = ["subnet-public-subnet"]
    service_account = google_service_account.cluster_gcp.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]
  }
  provider = "google.europe-west1"
}
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
  name                            = "example-vn"
  routing_mode                    = "REGIONAL"
  description                     = "Managed by Multy"
  auto_create_subnetworks         = false
  delete_default_routes_on_create = true
  provider                        = "google.europe-west1"
}
resource "aws_subnet" "public_subnet_aws-1" {
  tags                    = { "Name" = "public-subnet-1" }
  cidr_block              = "10.0.0.0/25"
  vpc_id                  = aws_vpc.example_vn_aws.id
  availability_zone       = "eu-west-1a"
  map_public_ip_on_launch = true
  provider                = "aws.eu-west-1"
}
resource "aws_subnet" "public_subnet_aws-2" {
  tags                    = { "Name" = "public-subnet-2" }
  cidr_block              = "10.0.0.128/26"
  vpc_id                  = aws_vpc.example_vn_aws.id
  availability_zone       = "eu-west-1b"
  map_public_ip_on_launch = true
  provider                = "aws.eu-west-1"
}
resource "aws_subnet" "public_subnet_aws-3" {
  tags                    = { "Name" = "public-subnet-3" }
  cidr_block              = "10.0.0.192/26"
  vpc_id                  = aws_vpc.example_vn_aws.id
  availability_zone       = "eu-west-1c"
  map_public_ip_on_launch = true
  provider                = "aws.eu-west-1"
}
resource "azurerm_subnet" "public_subnet_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "public-subnet"
  address_prefixes     = ["10.0.0.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "google_compute_subnetwork" "public_subnet_gcp" {
  name                     = "public-subnet"
  ip_cidr_range            = "10.0.0.0/24"
  network                  = google_compute_network.example_vn_gcp.id
  private_ip_google_access = true
  provider                 = "google.europe-west1"
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
resource "aws_route_table_association" "rta_aws-1" {
  subnet_id      = aws_subnet.public_subnet_aws-1.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "aws_route_table_association" "rta_aws-2" {
  subnet_id      = aws_subnet.public_subnet_aws-2.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "aws_route_table_association" "rta_aws-3" {
  subnet_id      = aws_subnet.public_subnet_aws-3.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "azurerm_subnet_route_table_association" "public_subnet_azure" {
  subnet_id      = azurerm_subnet.public_subnet_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
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
