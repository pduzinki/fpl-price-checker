package main

import (
	"github.com/pduzinki/fpl-price-checker/internal/di"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	fetchService := di.NewFetchServiceLambdas()

	lambda.Start(fetchService.Fetch)
}
