terraform {
  required_version = "~> 1.3"

  # cloud {
  #   organization = "fpc-org-1"
  #   workspaces {
  #     name = "fpc"
  #   }
  # }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.60.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.4.0"
    }
  }
}