resource "aws_iam_role" "cluster_aws" {
  tags               = { "Name" = "iam_for_k8cluster_cluster_aws" }
  name               = "iam_for_k8cluster_cluster_aws"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"eks.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_role_policy_attachment" "cluster_aws_AmazonEKSClusterPolicy" {
  role       = aws_iam_role.cluster_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}
resource "aws_iam_role_policy_attachment" "cluster_aws_AmazonEKSVPCResourceController" {
  role       = aws_iam_role.cluster_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
}
resource "aws_eks_cluster" "cluster_aws" {
  tags     = { "Name" = "cluster_aws" }
  role_arn = aws_iam_role.cluster_aws.arn
  vpc_config {
    subnet_ids = [aws_subnet.private_subnet_aws.id, aws_subnet.public_subnet_aws.id]
  }
  name = "cluster_aws"
}
resource "aws_iam_role" "node_pool_aws" {
  tags               = { "Name" = "iam_for_k8nodepool_node_pool_aws" }
  name               = "iam_for_k8nodepool_node_pool_aws"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_role_policy_attachment" "node_pool_aws_AmazonEKSWorkerNodePolicy" {
  role       = aws_iam_role.node_pool_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
}
resource "aws_iam_role_policy_attachment" "node_pool_aws_AmazonEKS_CNI_Policy" {
  role       = aws_iam_role.node_pool_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
}
resource "aws_iam_role_policy_attachment" "node_pool_aws_AmazonEC2ContainerRegistryReadOnly" {
  role       = aws_iam_role.node_pool_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}
resource "aws_eks_node_group" "node_pool_aws" {
  cluster_name    = aws_eks_cluster.cluster_aws.id
  node_group_name = "node_pool_aws"
  node_role_arn   = aws_iam_role.node_pool_aws.arn
  subnet_ids      = [
    aws_subnet.private_subnet_aws.id,
    aws_subnet.public_subnet_aws.id
  ]
  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }
  instance_types = ["t2.medium"]
}
resource "aws_vpc" "example_vn_aws" {
  tags                 = { "Name" = "example_vn" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags   = { "Name" = "example_vn" }
  vpc_id = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  tags   = { "Name" = "example_vn" }
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
resource "aws_route_table" "rt_aws" {
  tags   = { "Name" = "test-rt" }
  vpc_id = aws_vpc.example_vn_aws.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
}
resource "aws_route_table_association" "public_subnet_aws" {
  subnet_id      = aws_subnet.public_subnet_aws.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_subnet" "private_subnet_aws" {
  tags                    = { "Name" = "private-subnet" }
  cidr_block              = "10.0.1.0/24"
  vpc_id                  = aws_vpc.example_vn_aws.id
  availability_zone       = "eu-west-1a"
  map_public_ip_on_launch = true
}
resource "aws_subnet" "public_subnet_aws" {
  tags                    = { "Name" = "public-subnet" }
  cidr_block              = "10.0.2.0/24"
  vpc_id                  = aws_vpc.example_vn_aws.id
  availability_zone       = "eu-west-1b"
  map_public_ip_on_launch = true
}
resource "azurerm_kubernetes_cluster" "cluster_azure" {
  resource_group_name = azurerm_resource_group.ks-rg.name
  name                = "cluster_azure"
  location            = "northeurope"
  default_node_pool {
    name                = "node_pool_azure"
    node_count          = 1
    max_count           = 1
    min_count           = 1
    enable_auto_scaling = true
    vm_size             = "Standard_A2_v2"
  }
  dns_prefix = "cluster_azure"
  identity {
    type = "SystemAssigned"
  }
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
resource "azurerm_resource_group" "ks-rg" {
  name     = "ks-rg"
  location = "northeurope"
}
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "test-rt"
  location            = "northeurope"
  route {
    name           = "internet"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "Internet"
  }
}
resource "azurerm_subnet_route_table_association" "public_subnet_azure" {
  subnet_id      = azurerm_subnet.public_subnet_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "azurerm_subnet" "private_subnet_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "private-subnet"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "private_subnet_azure" {
  subnet_id      = azurerm_subnet.private_subnet_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
resource "azurerm_subnet" "public_subnet_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "public-subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_resource_group" "vn-rg" {
  name     = "vn-rg"
  location = "northeurope"
}
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {
  }
}