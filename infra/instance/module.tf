terraform {
  required_version = ">= 0.12"
}

provider "aws" {
  region = var.region
}

data "aws_availability_zones" "available" {
  state             = "available"
}

data "aws_ami" "default" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-2.0.2020*"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  owners = ["amazon"]
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "build/bin/app"
  output_path = "build/bin/app.zip"
}

resource "aws_instance" "server" {
  count                  = 1 
  instance_type          = var.instance_type
  ami                    = data.aws_ami.default.id

  credit_specification {
    cpu_credits = "standard"
  }

  tags = {
    Name = "cdn-server-${element(data.aws_availability_zones.available.names, count.index)}-${count.index}"
  }
}
