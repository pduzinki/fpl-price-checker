output "fpc_lambda_zip_storage_bucket_name" {
  description = "name of the s3 bucket used to store lambdas zips."

  value = random_pet.fpc_lambda_zip_storage_bucket_name
}