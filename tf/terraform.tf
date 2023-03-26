terraform {
  required_version = "~> 1.3"

  cloud {
    workspaces {
      name = "fpc"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.60.0"
    }
  }
}