provider "aws" {
  region = var.aws_region
}

# s3 storage for fpc data
resource "aws_s3_bucket" "fpc_bucket" {
  bucket = var.AWS_S3_BUCKET
}

# s3 storage for lambdas zips
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

# lambda fetch
resource "aws_s3_object" "lambda_fetch" {
  bucket = aws_s3_bucket.lambda_storage_bucket.id

  key    = "fetch.zip"
  source = "${path.module}/../build/lambdas/fetch.zip"

  # etag = filemd5("${path.module}/../build/lambdas/fetch.zip")`
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
  name                = "fetch-tf-cron"
  description         = "event123123"
  schedule_expression = "cron(30 3 * * ? *)"
}

resource "aws_cloudwatch_event_target" "lambda_target" {
  rule = aws_cloudwatch_event_rule.fetch-cron.name
  arn  = aws_lambda_function.fetch.arn
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_fetch_exec.name
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fetch.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.fetch-cron.arn
}

# lambda generate
resource "aws_s3_object" "lambda_generate" {
  bucket = aws_s3_bucket.lambda_storage_bucket.id

  key    = "generate.zip"
  source = "${path.module}/../build/lambdas/generate.zip"

  # etag = filemd5("${path.module}/../build/lambdas/generate.zip")
}

resource "aws_lambda_function" "generate" {
  function_name = "generate_tf"

  s3_bucket = aws_s3_bucket.lambda_storage_bucket.id
  s3_key    = aws_s3_object.lambda_generate.key

  runtime = "go1.x"
  handler = "generate"

  # source_code_hash = 

  role = aws_iam_role.lambda_generate_exec.arn

  environment {
    variables = {
      AWS_ID        = "${var.AWS_ID}"
      AWS_S3_BUCKET = "${var.AWS_S3_BUCKET}"
      AWS_SECRET    = "${var.AWS_SECRET}"
    }
  }
}

resource "aws_cloudwatch_log_group" "generate" {
  name = "/aws/lambda/${aws_lambda_function.generate.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "lambda_generate_exec" {
  name = "lambda_generate"

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

resource "aws_cloudwatch_event_rule" "generate-cron" {
  name                = "generate-tf-cron"
  description         = "event"
  schedule_expression = "cron(35 3 * * ? *)"
}

resource "aws_cloudwatch_event_target" "lambda_generate_target" {
  rule = aws_cloudwatch_event_rule.generate-cron.name
  arn  = aws_lambda_function.generate.arn
}

resource "aws_iam_role_policy_attachment" "lambda_policy_generate" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_generate_exec.name
}

resource "aws_lambda_permission" "allow_eventbridge_generate" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.generate.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.generate-cron.arn
}


# lambda get
resource "aws_s3_object" "lambda_get" {
  bucket = aws_s3_bucket.lambda_storage_bucket.id

  key    = "get.zip"
  source = "${path.module}/../build/lambdas/get.zip"

  # etag = filemd5("${path.module}/../build/lambdas/get.zip")
}

resource "aws_lambda_function" "get" {
  function_name = "get_tf"

  s3_bucket = aws_s3_bucket.lambda_storage_bucket.id
  s3_key    = aws_s3_object.lambda_get.key

  runtime = "go1.x"
  handler = "get"

  # source_code_hash = 

  role = aws_iam_role.lambda_get_exec.arn

  environment {
    variables = {
      AWS_ID        = "${var.AWS_ID}"
      AWS_S3_BUCKET = "${var.AWS_S3_BUCKET}"
      AWS_SECRET    = "${var.AWS_SECRET}"
    }
  }
}

resource "aws_cloudwatch_log_group" "get" {
  name = "/aws/lambda/${aws_lambda_function.get.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "lambda_get_exec" {
  name = "lambda_get"

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

resource "aws_iam_role_policy_attachment" "lambda_policy_get" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_get_exec.name
}

# api gateway
resource "aws_apigatewayv2_api" "lambda" {
  name          = "fpc_tf"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "lambda" {
  api_id = aws_apigatewayv2_api.lambda.id

  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )
  }
}

resource "aws_apigatewayv2_integration" "hello_world" {
  api_id = aws_apigatewayv2_api.lambda.id

  integration_uri    = aws_lambda_function.get.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"

  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "hello_world" {
  api_id = aws_apigatewayv2_api.lambda.id

  route_key = "GET /latest"
  target    = "integrations/${aws_apigatewayv2_integration.hello_world.id}"
}

resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.lambda.name}"

  retention_in_days = 30
}

resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}
