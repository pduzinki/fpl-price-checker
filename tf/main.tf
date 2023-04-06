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

  etag = filemd5("${path.module}/../build/lambdas/fetch.zip")
}

resource "aws_lambda_function" "fetch" {
  function_name = "fetch_tf"

  s3_bucket = aws_s3_bucket.lambda_storage_bucket.id
  s3_key    = aws_s3_object.lambda_fetch.key

  runtime = "go1.x"
  handler = "fetch"

  # source_code_hash = 

  role = aws_iam_role.lambda_fetch_exec.arn

  environment {
    variables = {
      AWS_ID        = "${var.AWS_ID}"
      AWS_S3_BUCKET = "${var.AWS_S3_BUCKET}"
      AWS_SECRET    = "${var.AWS_SECRET}"
    }
  }
}

resource "aws_cloudwatch_log_group" "fetch" {
  name = "/aws/lambda/${aws_lambda_function.fetch.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "lambda_fetch_exec" {
  name = "lambda_fetch"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      }
    ]
  })
}