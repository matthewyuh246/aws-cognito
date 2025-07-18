provider "aws" {
    region = "ap-northeast-1"
    profile = "admin"
}

# data
data "aws_region" "current" {}

# variables
variable "project" {
    type = string
}

variable "environment" {
    type = string
}

#　ローカル変数(= 再利用可能な値)
locals {
  vpc_cidr_block             = "10.0.0.0/16"
  public_subnet_cidr_blocks  = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnet_cidr_blocks = ["10.0.3.0/24", "10.0.4.0/24"]
  availability_zones         = ["ap-northeast-1a", "ap-northeast-1c"]
}