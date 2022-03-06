resource "aws_iam_role" "example_aws" {
  tags               = { "Name" = "iam_for_k8cluster_example" }
  name               = "iam_for_k8cluster_example"
  assume_role_policy = "{\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"eks.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_role_policy_attachment" "example_aws_AmazonEKSClusterPolicy" {
  role       = aws_iam_role.example_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}
resource "aws_iam_role_policy_attachment" "example_aws_AmazonEKSVPCResourceController" {
  role       = aws_iam_role.example_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
}
resource "aws_eks_cluster" "example" {
  tags     = { "Name" = "example" }
  role_arn = aws_iam_role.example_aws.arn
  vpc_config {
    subnet_ids = ["${aws_subnet.subnet1_aws.id}", "${aws_subnet.subnet2_aws.id}"]
  }
  name = "example"
}
resource "aws_iam_role" "example_pool_aws" {
  tags               = { "Name" = "iam_for_k8nodepool_example" }
  name               = "iam_for_k8nodepool_example"
  assume_role_policy = "{\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_role_policy_attachment" "example_pool_aws_AmazonEKSWorkerNodePolicy" {
  role       = aws_iam_role.example_pool_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
}
resource "aws_iam_role_policy_attachment" "example_pool_aws_AmazonEKS_CNI_Policy" {
  role       = aws_iam_role.example_pool_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
}
resource "aws_iam_role_policy_attachment" "example_pool_aws_AmazonEC2ContainerRegistryReadOnly" {
  role       = aws_iam_role.example_pool_aws.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}
resource "aws_eks_node_group" "example_pool" {
  cluster_name    = "example"
  node_group_name = "example"
  node_role_arn   = aws_iam_role.example_pool_aws.arn
  subnet_ids      = ["${aws_subnet.subnet1_aws.id}", "${aws_subnet.subnet2_aws.id}"]
  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }
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
  vpc_id = "${aws_vpc.example_vn_aws.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
}
resource "aws_route_table_association" "rta_aws" {
  subnet_id      = "${aws_subnet.subnet2_aws.id}"
  route_table_id = "${aws_route_table.rt_aws.id}"
}
resource "aws_subnet" "subnet1_aws" {
  tags                    = { "Name" = "private-subnet" }
  cidr_block              = "10.0.1.0/24"
  vpc_id                  = aws_vpc.example_vn_aws.id
  availability_zone       = "eu-west-1a"
  map_public_ip_on_launch = true
}
resource "aws_subnet" "subnet2_aws" {
  tags                    = { "Name" = "public-subnet" }
  cidr_block              = "10.0.2.0/24"
  vpc_id                  = aws_vpc.example_vn_aws.id
  availability_zone       = "eu-west-1b"
  map_public_ip_on_launch = true
}
provider "aws" {
  region = "eu-west-1"
}
