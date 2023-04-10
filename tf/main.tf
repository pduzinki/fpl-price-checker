provider "aws" {
  region = var.aws_region
}

# s3 storage for fpc data
resource "aws_s3_bucket" "fpc_data_bucket" {
  bucket = var.AWS_S3_BUCKET
}

# s3 storage for fpc lambdas zips
resource "random_pet" "fpc_lambda_zip_storage_bucket_name" {
  prefix = "fpc-lambda-zip-storage"
}

resource "aws_s3_bucket" "fpc_lambda_zip_storage_bucket" {
  bucket = random_pet.fpc_lambda_zip_storage_bucket_name.id
}

resource "aws_s3_bucket_acl" "fpc_lambda_zip_storage_bucket_acl" {
  bucket = aws_s3_bucket.fpc_lambda_zip_storage_bucket.id
  acl    = "private"
}

# lambda for fpc 'fetch' command
resource "aws_s3_object" "fpc_fetch" {
  bucket = aws_s3_bucket.fpc_lambda_zip_storage_bucket.id

  key    = "fetch.zip"
  source = "${path.module}/../build/lambdas/fetch.zip"
}

resource "aws_lambda_function" "fpc_fetch" {
  function_name = "fpc-fetch"

  s3_bucket = aws_s3_bucket.fpc_lambda_zip_storage_bucket.id
  s3_key    = aws_s3_object.fpc_fetch.key

  runtime = "go1.x"
  handler = "fetch"

  role = aws_iam_role.fpc_fetch.arn

  environment {
    variables = {
      AWS_ID        = "${var.AWS_ID}"
      AWS_S3_BUCKET = "${var.AWS_S3_BUCKET}"
      AWS_SECRET    = "${var.AWS_SECRET}"
    }
  }
}

resource "aws_cloudwatch_log_group" "fpc_fetch" {
  name = "/aws/lambda/${aws_lambda_function.fpc_fetch.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "fpc_fetch" {
  name = "fpc-fetch"

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

resource "aws_cloudwatch_event_rule" "fpc_fetch" {
  name                = "fpc-fetch-eventbridge-rule"
  description         = "eventbridge schedule rule for fpc 'fetch' command"
  schedule_expression = "cron(30 3 * * ? *)"
}

resource "aws_cloudwatch_event_target" "fpc_fetch" {
  rule = aws_cloudwatch_event_rule.fpc_fetch.name
  arn  = aws_lambda_function.fpc_fetch.arn
}

resource "aws_iam_role_policy_attachment" "fpc_fetch_lambda_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.fpc_fetch.name
}

resource "aws_lambda_permission" "fpc_fetch_allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fpc_fetch.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.fpc_fetch.arn
}

# lambda for fpc 'generate' command
resource "aws_s3_object" "fpc_generate" {
  bucket = aws_s3_bucket.fpc_lambda_zip_storage_bucket.id

  key    = "generate.zip"
  source = "${path.module}/../build/lambdas/generate.zip"
}

resource "aws_lambda_function" "fpc_generate" {
  function_name = "fpc-generate"

  s3_bucket = aws_s3_bucket.fpc_lambda_zip_storage_bucket.id
  s3_key    = aws_s3_object.fpc_generate.key

  runtime = "go1.x"
  handler = "generate"

  role = aws_iam_role.fpc_generate.arn

  environment {
    variables = {
      AWS_ID        = "${var.AWS_ID}"
      AWS_S3_BUCKET = "${var.AWS_S3_BUCKET}"
      AWS_SECRET    = "${var.AWS_SECRET}"
    }
  }
}

resource "aws_cloudwatch_log_group" "fpc_generate" {
  name = "/aws/lambda/${aws_lambda_function.fpc_generate.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "fpc_generate" {
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

resource "aws_cloudwatch_event_rule" "fpc_generate" {
  name                = "fpc-generate-eventbridge-rule"
  description         = "eventbridge schedule rule for fpc 'generate' command"
  schedule_expression = "cron(35 3 * * ? *)"
}

resource "aws_cloudwatch_event_target" "fpc_generate" {
  rule = aws_cloudwatch_event_rule.fpc_generate.name
  arn  = aws_lambda_function.fpc_generate.arn
}

resource "aws_iam_role_policy_attachment" "fpc_generate_lambda_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.fpc_generate.name
}

resource "aws_lambda_permission" "fpc_getch_allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fpc_generate.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.fpc_generate.arn
}


# lambda for fpc 'get' command
resource "aws_s3_object" "fpc_get" {
  bucket = aws_s3_bucket.fpc_lambda_zip_storage_bucket.id

  key    = "get.zip"
  source = "${path.module}/../build/lambdas/get.zip"
}

resource "aws_lambda_function" "fpc_get" {
  function_name = "fpc-get"

  s3_bucket = aws_s3_bucket.fpc_lambda_zip_storage_bucket.id
  s3_key    = aws_s3_object.fpc_get.key

  runtime = "go1.x"
  handler = "get"

  role = aws_iam_role.fpc_get.arn

  environment {
    variables = {
      AWS_ID        = "${var.AWS_ID}"
      AWS_S3_BUCKET = "${var.AWS_S3_BUCKET}"
      AWS_SECRET    = "${var.AWS_SECRET}"
    }
  }
}

resource "aws_cloudwatch_log_group" "fpc_get" {
  name = "/aws/lambda/${aws_lambda_function.fpc_get.function_name}"

  retention_in_days = 30
}

resource "aws_iam_role" "fpc_get" {
  name = "fpc-get"

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

resource "aws_iam_role_policy_attachment" "fpc_get_lambda_policyt" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.fpc_get.name
}

resource "aws_lambda_permission" "fpc_get_allow_api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fpc_get.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.fpc_api_gateway.execution_arn}/*/*"
}

# api gateway to access reports via fpc 'get' lambda
resource "aws_apigatewayv2_api" "fpc_api_gateway" {
  name          = "fpc-api-gateway"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "fpc_api_gateway_stage" {
  api_id = aws_apigatewayv2_api.fpc_api_gateway.id

  name        = "prod"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.fpc_api_gateway.arn

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

resource "aws_apigatewayv2_integration" "fpc_api_gateway" {
  api_id = aws_apigatewayv2_api.fpc_api_gateway.id

  integration_uri    = aws_lambda_function.fpc_get.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"

  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "fpc_get_latest" {
  api_id = aws_apigatewayv2_api.fpc_api_gateway.id

  route_key = "GET /latest"
  target    = "integrations/${aws_apigatewayv2_integration.fpc_api_gateway.id}"
}

resource "aws_cloudwatch_log_group" "fpc_api_gateway" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.fpc_api_gateway.name}"

  retention_in_days = 30
}
