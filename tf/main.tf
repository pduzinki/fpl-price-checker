provider "aws" {
  region = var.aws_region
}

resource "random_pet" "lambda_storage_bucket_name" {
  prefix = "fpc-lambda-storage"
}

resource "aws_s3_bucket" "lambda_storage_bucket" {
  bucket = random_pet.lambda_storage_bucket_name.id
}

resource "aws_s3_bucket_acl" "lambda_storage_bucket_acl" {
  bucket = aws_s3_bucket.lambda_storage_bucket.id
  acl    = "private"
}

resource "aws_s3_object" "lambda_fetch" {
  bucket = aws_s3_bucket.lambda_storage_bucket.id

  key    = "fetch.zip"
  source = "${path.module}/../build/lambdas/fetch.zip"
}