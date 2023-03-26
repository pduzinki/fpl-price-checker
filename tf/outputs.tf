output "lambda-storage-bucket-name" {
  description = "name of the s3 bucket used to store lambdas code."

  value = random_pet.lambda_storage_bucket_name
}