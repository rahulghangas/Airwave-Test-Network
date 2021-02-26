module "cdn-us-east-1" {
  source        = "./instance"
  region        = "us-east-1"
  instance_type = "t2.micro"
}

module "cdn-us-east-2" {
  source = "./instance"
  region = "us-east-2"
}

module "cdn-us-west-1" {
  source = "./instance"
  region = "us-west-1"
}

module "cdn-us-west-2" {
  source = "./instance"
  region = "us-west-2"
}

module "cdn-eu-west-1" {
  source = "./instance"
  region = "eu-west-1"
}

module "cdn-eu-central-1" {
  source = "./instance"
  region = "eu-central-1"
}

module "cdn-eu-north-1" {
  source = "./instance"
  region = "eu-north-1"
}

module "cdn-ap-southeast-1" {
  source = "./instance"
  region = "ap-southeast-1"
}

module "cdn-sa-east-1" {
  source = "./instance"
  region = "sa-east-1"
}
