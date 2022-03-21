variable bucket_name {
  type    = string
  default = "multy-users-tfstate"
}
resource "aws_vpc" "main_vpc" {
  tags                 = { "Name" = "backend" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags   = { "Name" = "backend" }
  vpc_id = aws_vpc.main_vpc.id
}
resource "aws_default_security_group" "example_vn_aws" {
  tags   = { "Name" = "backend" }
  vpc_id = aws_vpc.main_vpc.id
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
resource "aws_security_group" "nsg2_aws" {
  tags   = { "Name" = "backend" }
  vpc_id = aws_vpc.main_vpc.id
  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    protocol    = "tcp"
    from_port   = 8000
    to_port     = 8000
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    protocol    = "tcp"
    from_port   = 443
    to_port     = 443
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["10.0.0.0/16"]
  }
  egress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "tcp"
    from_port   = 8000
    to_port     = 8000
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "tcp"
    from_port   = 443
    to_port     = 443
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["10.0.0.0/16"]
  }
}
resource "aws_route_table" "rt_aws" {
  tags   = { "Name" = "backend" }
  vpc_id = aws_vpc.main_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
}
resource "aws_route_table_association" "rta" {
  subnet_id      = aws_subnet.subnet.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_subnet" "subnet" {
  tags       = { "Name" = "backend" }
  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.main_vpc.id
}
resource "aws_iam_instance_profile" "iam_instance_profile" {
  name = "backend_iam_profile"
  role = aws_iam_role.vm_iam.name
}
data "aws_caller_identity" "vm_aws" {}
data "aws_region" "vm_aws" {}
resource "aws_iam_role" "vm_iam" {
  tags               = { "Name" = "backend" }
  name               = "backend_iam"
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17"
    "Statement" : [
      {
        "Action" : ["sts:AssumeRole"],
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "ec2.amazonaws.com"
        }
      }
    ],
  })
  inline_policy {
    name   = "vault_policy"
    policy = jsonencode({
      "Version" : "2012-10-17"
      "Statement" : [
        {
          "Effect" : "Allow",
          "Action" : "s3:*",
          "Resource" : aws_s3_bucket.tfstate_bucket.arn
        },
        {
          "Effect" : "Allow",
          "Action" : [
            "s3:ListAccessPointsForObjectLambda",
            "s3:GetAccessPoint",
            "s3:PutAccountPublicAccessBlock",
            "s3:ListAccessPoints",
            "dynamodb:ListTables",
            "s3:ListJobs",
            "dynamodb:ListBackups",
            "s3:PutStorageLensConfiguration",
            "dynamodb:PurchaseReservedCapacityOfferings",
            "s3:ListMultiRegionAccessPoints",
            "dynamodb:ListStreams",
            "s3:ListStorageLensConfigurations",
            "dynamodb:ListContributorInsights",
            "s3:GetAccountPublicAccessBlock",
            "dynamodb:DescribeReservedCapacityOfferings",
            "s3:ListAllMyBuckets",
            "dynamodb:ListGlobalTables",
            "s3:PutAccessPointPublicAccessBlock",
            "dynamodb:DescribeReservedCapacity",
            "s3:CreateJob",
            "dynamodb:DescribeLimits",
            "dynamodb:ListExports"
          ],
          "Resource" : "*"
        },
        {
          "Effect" : "Allow",
          "Action" : "dynamodb:*",
          "Resource" : aws_dynamodb_table.user_ddb.arn
        },
      ],
    })
  }
}
resource "aws_key_pair" "vm" {
  tags       = { "Name" = "backend" }
  key_name   = "vm_multy"
  public_key = file("./ssh_key.pub")
}
#resource "aws_instance" "vm" {
#  tags             = { "Name" = "backend" }
#  ami              = "ami-0015a39e4b7c0966f" # Ubuntu Server 20.04 LTS (HVM), SSD Volume Type
#  instance_type    = "t2.nano"
#  subnet_id        = aws_subnet.subnet.id
#  user_data_base64 = base64encode(templatefile("init.sh", {
#    "s3_bucket_name" = var.bucket_name
#  }))
#  key_name             = aws_key_pair.vm.key_name
#  iam_instance_profile = aws_iam_instance_profile.iam_instance_profile.id
#}
resource "aws_s3_bucket" "tfstate_bucket" {
  tags   = { "Name" = "backend" }
  bucket = var.bucket_name
}
resource "aws_dynamodb_table" "user_ddb" {
  tags         = { "Name" = "backend" }
  name         = "user_table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "user_id"

  attribute {
    name = "user_id"
    type = "S"
  }
}
#resource "aws_eip" "ip_aws" {
#  tags     = { "Name" = "backend" }
#  instance = aws_instance.vm.id
#}
terraform {
  backend "s3" {
    bucket = "multy-tfstate"
    key    = "terraform.tfstate"
    region = "eu-west-2"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}
provider "aws" {
  region = "eu-west-2"
}
#output "aws_endpoint" {
#  value = aws_eip.ip_aws.public_ip
#}