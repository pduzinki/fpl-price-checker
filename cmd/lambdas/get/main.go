package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pduzinki/fpl-price-checker/pkg/di"
)

func main() {
	gs := di.NewGetServiceS3()

	lambda.Start(gs.GetLatestReport)
}
