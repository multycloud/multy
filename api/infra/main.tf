variable bucket_name {
  type    = string
  default = "multy-users-tfstate"
}
variable api_endpoint {
  type    = string
  default = "api2.multy.dev"
}
resource "aws_vpc" "main_vpc" {
  tags                 = { "Name" = "backend" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "vn_igw" {
  tags   = { "Name" = "backend" }
  vpc_id = aws_vpc.main_vpc.id
}
resource "aws_default_security_group" "vn_nsg" {
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
resource "aws_security_group" "nsg" {
  tags        = { "Name" = "backend" }
  vpc_id      = aws_vpc.main_vpc.id
  name        = "backend"
  description = "Managed by Multy"
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
  ingress {
    protocol    = "tcp"
    from_port   = 3306
    to_port     = 3306
    cidr_blocks = [aws_subnet.public_subnet.cidr_block]
  }
  egress {
    protocol    = "tcp"
    from_port   = 3306
    to_port     = 3306
    cidr_blocks = ["0.0.0.0/0"]
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
resource "aws_route_table" "public_rt" {
  tags   = { "Name" = "backend" }
  vpc_id = aws_vpc.main_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vn_igw.id
  }
}
resource "aws_route_table_association" "rta" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_rt.id
}
// todo: remove to put db private
resource "aws_route_table_association" "rta2" {
  subnet_id      = aws_subnet.private_subnet.id
  route_table_id = aws_route_table.public_rt.id
}
resource "aws_route_table_association" "rta3" {
  subnet_id      = aws_subnet.private_subnet2.id
  route_table_id = aws_route_table.public_rt.id
}
resource "aws_subnet" "public_subnet" {
  tags              = { "Name" = "backend" }
  cidr_block        = "10.0.1.0/24"
  availability_zone = "eu-west-2a"
  vpc_id            = aws_vpc.main_vpc.id
}
resource "aws_subnet" "private_subnet" {
  tags              = { "Name" = "backend" }
  cidr_block        = "10.0.2.0/24"
  availability_zone = "eu-west-2a"
  vpc_id            = aws_vpc.main_vpc.id
}
resource "aws_subnet" "private_subnet2" {
  tags              = { "Name" = "backend" }
  cidr_block        = "10.0.3.0/24"
  availability_zone = "eu-west-2b"
  vpc_id            = aws_vpc.main_vpc.id
}
resource "aws_iam_instance_profile" "iam_instance_profile" {
  name = "backend_iam_profile"
  role = aws_iam_role.vm_iam.name
}
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
          "Resource" : "${aws_s3_bucket.tfstate_bucket.arn}/*"
        },
        {
          "Effect" : "Allow",
          "Action" : [
            "s3:ListBucket"
          ],
          "Resource" : [
            "arn:aws:s3:::multy-internal"
          ]
        },
        {
          "Effect" : "Allow",
          "Action" : ["s3:Get*", "s3:List*",],
          "Resource" : "arn:aws:s3:::multy-internal/certs/*"
        },
        {
          "Effect" : "Allow",
          "Action" : [
            "s3:ListAccessPointsForObjectLambda",
            "s3:GetAccessPoint",
            "s3:PutAccountPublicAccessBlock",
            "s3:ListAccessPoints",
            "s3:ListJobs",
            "s3:PutStorageLensConfiguration",
            "s3:ListMultiRegionAccessPoints",
            "s3:ListStorageLensConfigurations",
            "s3:GetAccountPublicAccessBlock",
            "s3:ListAllMyBuckets",
            "s3:PutAccessPointPublicAccessBlock",
            "s3:CreateJob",
          ],
          "Resource" : "*"
        },
      ],
    })
  }
}
resource "aws_s3_bucket" "tfstate_bucket" {
  tags   = { "Name" = "backend" }
  bucket = var.bucket_name
}
resource "aws_key_pair" "vm" {
  tags       = { "Name" = "backend" }
  key_name   = "infra_key"
  public_key = file("./infra_key.pub")
}
resource "random_password" "password" {
  length = 16
}
resource "aws_ssm_parameter" "db_password" {
  name  = "/dev-multy/db-password"
  type  = "SecureString"
  value = random_password.password.result
}
resource "aws_instance" "vm" {
  tags             = { "Name" = "backend" }
  ami              = "ami-0015a39e4b7c0966f" # Ubuntu Server 20.04 LTS (HVM), SSD Volume Type
  instance_type    = "t2.nano"
  subnet_id        = aws_subnet.public_subnet.id
  user_data_base64 = base64encode(templatefile("init.sh", {
    "s3_bucket_name" = var.bucket_name
    "db_connection"  = "${aws_db_instance.db.username}:${random_password.password.result}@tcp(${aws_db_instance.db.address}:${aws_db_instance.db.port}/${aws_db_instance.db.name}"
    "api_endpoint"   = var.api_endpoint
  }))
  key_name             = aws_key_pair.vm.key_name
  iam_instance_profile = aws_iam_instance_profile.iam_instance_profile.id
}
resource "aws_db_subnet_group" "db_subnet_group" {
  tags = {
    "Name" = "example-db"
  }

  name = "example-db"

  subnet_ids = [
    aws_subnet.private_subnet.id,
    aws_subnet.private_subnet2.id,
  ]
}
resource "aws_db_instance" "db" {
  tags = {
    "Name" = "exampledb"
  }

  allocated_storage    = 10
  db_name              = "multydb"
  engine               = "mysql"
  engine_version       = "5.7"
  username             = "multyadmin"
  password             = random_password.password.result
  instance_class       = "db.t2.micro"
  identifier           = "multydb"
  skip_final_snapshot  = true
  db_subnet_group_name = aws_db_subnet_group.db_subnet_group.name
  publicly_accessible  = true
  deletion_protection  = true
}

resource "aws_eip" "vm_ip" {
  tags     = { "Name" = "backend" }
  instance = aws_instance.vm.id
}
data "aws_route53_zone" "primary" {
  name = "multy.dev"
}
resource "aws_route53_record" "server1-record" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = var.api_endpoint
  type    = "A"
  ttl     = "300"
  records = [aws_eip.vm_ip.public_ip]
}
resource "null_resource" "db_init" {
  triggers = {
    change = filemd5("../../db/init.sql")
  }
  provisioner "local-exec" {
    command = "mysql -h ${aws_db_instance.db.address} -P ${aws_db_instance.db.port} -u ${aws_db_instance.db.username} --password='${random_password.password.result}' -e 'source ../../db/init.sql'"
  }
}
terraform {
  backend "s3" {
    bucket         = "multy-tfstate"
    key            = "terraform.tfstate"
    region         = "eu-west-2"
    dynamodb_table = "terraform-lock"
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