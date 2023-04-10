variable "aws_region" {
  type    = string
  default = "eu-north-1"
}

variable "AWS_ID" {
  type      = string
  sensitive = true
  default   = "TODO define me in terraform.tfvars"
}

variable "AWS_S3_BUCKET" {
  type      = string
  sensitive = true
  default   = "TODO define me in terraform.tfvars"
}

variable "AWS_SECRET" {
  type      = string
  sensitive = true
  default   = "TODO define me in terraform.tfvars"
}