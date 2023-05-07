package main

import (
	"context"
	"fmt"

	"github.com/pduzinki/fpl-price-checker/internal/di"
	"github.com/pduzinki/fpl-price-checker/internal/domain"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

func main() {
	gs := di.NewGetServiceLambdas()

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (domain.PriceChangeReport, error) {
		log.Info().Msg(fmt.Sprintf("api gateway request: %v", req))

		if date, prs := req.PathParameters["date"]; prs {
			return gs.GetReportByDate(ctx, date)
		}

		return gs.GetLatestReport(ctx)
	})
}
