package main

import (
	"github.com/pduzinki/fpl-price-checker/internal/di"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	gs := di.NewGenerateServiceLambdas()

	lambda.Start(gs.GeneratePriceReport)
}
