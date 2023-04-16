output "fpc_lambda_zip_storage_bucket_name" {
  description = "name of the s3 bucket used to store lambdas zips."

  value = random_pet.fpc_lambda_zip_storage_bucket_name.id
}

output "fpc_api_gateway_get_latest_report" {
  description = "url of the fpc api gateway /latest endpoint"

  value =format("%s%s", aws_apigatewayv2_stage.fpc_api_gateway_stage.invoke_url, split(" ", aws_apigatewayv2_route.fpc_get_latest.route_key)[1])
}