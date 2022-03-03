#resource "aws_eks_cluster" "example" {
#  name     = "example"
#  role_arn = aws_iam_role.example.arn
#
#  vpc_config {
#    subnet_ids = [aws_subnet.subnet1_aws.id, aws_subnet.subnet2_aws.id]
#  }
#
#  # Ensure that IAM Role permissions are created before and deleted after EKS Cluster handling.
#  # Otherwise, EKS will not be able to properly delete EKS managed EC2 infrastructure such as Security Groups.
#  depends_on = [
#    aws_iam_role_policy_attachment.example-AmazonEKSClusterPolicy,
#    aws_iam_role_policy_attachment.example-AmazonEKSVPCResourceController,
#    aws_cloudwatch_log_group.example
#  ]
#  enabled_cluster_log_types = ["api", "audit"]
#}
#
#resource "aws_vpc" "example_vn_aws" {
#  tags = {
#    "Name" = "example_vn"
#  }
#
#  cidr_block           = "10.0.0.0/16"
#  enable_dns_hostnames = true
#}
#
#resource "aws_subnet" "subnet1_aws" {
#  tags = {
#    "Name" = "private-subnet1"
#  }
#
#  cidr_block        = "10.0.1.0/24"
#  vpc_id            = aws_vpc.example_vn_aws.id
#  availability_zone = "eu-west-1a"
#}
#resource "aws_subnet" "subnet2_aws" {
#  tags = {
#    "Name" = "public-subnet2"
#  }
#
#  cidr_block        = "10.0.2.0/24"
#  vpc_id            = aws_vpc.example_vn_aws.id
#  availability_zone = "eu-west-1b"
#  map_public_ip_on_launch = tre
#}
#
#resource "aws_route_table_association" "rta_aws" {
#  subnet_id      = "${aws_subnet.subnet2_aws.id}"
#  route_table_id = "${aws_route_table.rt_aws.id}"
#}
#
#
#resource "aws_route_table" "rt_aws" {
#  tags = {
#    "Name" = "test-rt"
#  }
#
#  vpc_id = "${aws_vpc.example_vn_aws.id}"
#
#  route {
#    cidr_block = "0.0.0.0/0"
#    gateway_id = aws_internet_gateway.example_vn_aws.id
#  }
#}
#
#resource "aws_internet_gateway" "example_vn_aws" {
#  tags = {
#    "Name" = "example_vn"
#  }
#
#  vpc_id = aws_vpc.example_vn_aws.id
#}
#
#output "endpoint" {
#  value = aws_eks_cluster.example.endpoint
#}
#
#output "kubeconfig-certificate-authority-data" {
#  value = aws_eks_cluster.example.certificate_authority[0].data
#}
#resource "aws_iam_role" "example" {
#  name = "eks-node-group-example"
#
#  assume_role_policy = jsonencode({
#    Statement = [{
#      Action = "sts:AssumeRole"
#      Effect = "Allow"
#      Principal = {
#        Service = "ec2.amazonaws.com"
#      }
#    },{
#      Action = "sts:AssumeRole"
#      Effect = "Allow"
#      Principal = {
#        Service = "eks.amazonaws.com"
#      }
#    }]
#    Version = "2012-10-17"
#  })
#}
#
#resource "aws_iam_role_policy_attachment" "example-AmazonEKSClusterPolicy" {
#  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
#  role       = aws_iam_role.example.name
#}
#
## Optionally, enable Security Groups for Pods
## Reference: https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html
#resource "aws_iam_role_policy_attachment" "example-AmazonEKSVPCResourceController" {
#  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
#  role       = aws_iam_role.example.name
#}
#
#resource "aws_iam_role_policy_attachment" "example-AmazonEKSWorkerNodePolicy" {
#  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
#  role       = aws_iam_role.example.name
#}
#
#resource "aws_iam_role_policy_attachment" "example-AmazonEKS_CNI_Policy" {
#  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
#  role       = aws_iam_role.example.name
#}
#
#resource "aws_iam_role_policy_attachment" "example-AmazonEC2ContainerRegistryReadOnly" {
#  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
#  role       = aws_iam_role.example.name
#}
#
#provider "aws" {
#  region = "eu-west-1"
#}
#
#resource "aws_cloudwatch_log_group" "example" {
#  # The log group name format is /aws/eks/<cluster-name>/cluster
#  # Reference: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
#  name              = "/aws/eks/example/cluster"
#  retention_in_days = 7
#
#  # ... potentially other configuration ...
#}
#
#resource "aws_eks_node_group" "example" {
#  cluster_name    = aws_eks_cluster.example.name
#  node_group_name = "example"
#  node_role_arn   = aws_iam_role.example.arn
#  subnet_ids = [aws_subnet.subnet1_aws.id, aws_subnet.subnet2_aws.id]
#
#  scaling_config {
#    desired_size = 1
#    max_size     = 1
#    min_size     = 1
#  }
#
#  update_config {
#    max_unavailable = 1
#  }
#
#  depends_on = [
#    aws_iam_role_policy_attachment.example-AmazonEKSWorkerNodePolicy,
#    aws_iam_role_policy_attachment.example-AmazonEKS_CNI_Policy,
#    aws_iam_role_policy_attachment.example-AmazonEC2ContainerRegistryReadOnly,
#  ]
#}