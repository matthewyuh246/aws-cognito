# vpc
resource "aws_vpc" "main" {
  cidr_block           = local.vpc_cidr_block
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "${var.project}-${var.environment}-vpc"
  }
}

# subnet
resource "aws_subnet" "public" {
  count             = length(local.public_subnet_cidr_blocks)
  vpc_id            = aws_vpc.main.id
  cidr_block        = element(local.public_subnet_cidr_blocks, count.index)
  availability_zone = element(local.availability_zones, count.index)

  tags = {
    Name = "${var.project}-${var.environment}-public-subnet0${count.index}"
  }
}

resource "aws_subnet" "private" {
  count             = length(local.private_subnet_cidr_blocks)
  vpc_id            = aws_vpc.main.id
  cidr_block        = element(local.private_subnet_cidr_blocks, count.index)
  availability_zone = element(local.availability_zones, count.index)

  tags = {
    Name = "${var.project}-${var.environment}-private-subnet-${count.index}"
  }
}