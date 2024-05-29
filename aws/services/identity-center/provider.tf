terraform {
    backend "s3" {
        bucket = "the0x"
        key    = "services/identity-center/terraform.tfstate"
        region = "us-east-1"
    }

    required_providers {
        aws = {
            source  = "hashicorp/aws"
            version = "~> 5.0"
        }

        auth0 = {
            source = "auth0/auth0"
            version = "1.2.0"
        }
    }
}

provider "auth0" {
    domain = "https://dev-xrdsmokq.us.auth0.com/api/v2/"
}

provider "aws" {
    region = "us-east-1"
}