package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pduzinki/fpl-price-checker/pkg/di"
)

func main() {
	gs := di.NewGenerateServiceS3()

	lambda.Start(gs.GeneratePriceReport)
}
