provider "aws" {
  region = var.aws_region
}

resource "aws_s3_bucket" "fpc_bucket" {
  bucket = "${var.AWS_S3_BUCKET}"
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

resource "aws_cloudwatch_event_rule" "fetch-cron" {
  name        = "fetch-tf-cron"
  description = "event123123"
  schedule_expression = "cron(30 3 * * ? *)"
}

resource "aws_cloudwatch_event_target" "lambda_target" {
  rule      = aws_cloudwatch_event_rule.fetch-cron.name
  arn       = aws_lambda_function.fetch.arn
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role = aws_iam_role.lambda_fetch_exec.name
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fetch.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.fetch-cron.arn
}