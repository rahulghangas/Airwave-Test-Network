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

resource "aws_key_pair" "airwave" {
  key_name   = "airwave"
  public_key = file("key.pub")
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "../build/bin/app"
  output_path = "../build/bin/app.zip"
}

resource "aws_security_group" "airwave" {

  // Allow ssh connection through port 22
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // Listen on port 8080
  ingress {
    description = "TCP"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

}

resource "aws_instance" "airwave" {
  // Define number of instances to deploy in each region
  count                  = 1 
  instance_type          = var.instance_type
  key_name               = aws_key_pair.airwave.key_name
  ami                    = data.aws_ami.default.id

  // Expose to outside world
  vpc_security_group_ids = [
    aws_security_group.airwave.id
  ]

  // Set up ability to create a ssh connection
  connection {
    type        = "ssh"
    user        = "airwave"
    private_key = file("key")
    host        = self.public_ip
  }

  credit_specification {
    cpu_credits = "standard"
  }

  tags = {
    Name = "cdn-server-${element(data.aws_availability_zones.available.names, count.index)}-${count.index}"
  }
}

resource "aws_eip" "airwave" {
  // For some reason * does not seem to work, so used for_each 
  // to create multiple resources (corresponding to count)
  for_each   = { for idx, instance in aws_instance.airwave: idx => instance}
  vpc      = true
  instance = each.value.id
}
