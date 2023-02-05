package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pduzinki/fpl-price-checker/pkg/di"
)

func main() {
	fetchService := di.NewFetchServiceS3()

	lambda.Start(fetchService.Fetch)
}
