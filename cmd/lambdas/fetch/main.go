package main

import (
	"github.com/pduzinki/fpl-price-checker/pkg/di"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	fetchService := di.NewFetchServiceS3()

	lambda.Start(fetchService.Fetch)
}
