locals {
    vpc_cidr_blocks = [
        "10.0.0.0/20",
        "10.1.0.0/20",
        "10.2.0.0/20"
    ]
}

resource "aws_vpc" "vpc_config" {
    for_each = toset(local.vpc_cidr_blocks)
    cidr_block = each.value
}

resource "aws_subnet" "public_subnet_config" {
    for_each           = aws_vpc.vpc_config
    vpc_id             = each.value.id
    cidr_block         = cidrsubnet(each.value.cidr_block, 4, 0)
    availability_zone  = "us-east-1a"
}


resource "aws_subnet" "private_subnet_config" {
    for_each           = aws_vpc.vpc_config
    vpc_id             = each.value.id
    cidr_block         = cidrsubnet(each.value.cidr_block, 4, 1)
    availability_zone  = "us-east-1a"
}

resource "aws_ec2_transit_gateway" "transit_gateway" {
    description = "Transit Gateway"
}

resource "aws_ram_resource_share" "transit_gateway" {
  provider = aws.first
  name = "terraform-example"
  tags = {
    Name = "terraform-example"
  }
}

resource "aws_ram_resource_association" "example" {
  provider = aws.first
  resource_arn       = aws_ec2_transit_gateway.transit_gateway.arn
  resource_share_arn = aws_ram_resource_share.transit_gateway.id
}

resource "aws_ram_principal_association" "example" {
  provider = aws.first
  principal          = data.aws_caller_identity.second.account_id
  resource_share_arn = aws_ram_resource_share.transit_gateway.id
}
