variable "region" {
}

variable "servers_per_region" {
  default = 1
}

variable "instance_type" {
  default = "t3.nano"
}

variable "blacklisted_az" {
  default = ["us-west-1a", "us-east-1c", "ap-northeast-1a"]
}
