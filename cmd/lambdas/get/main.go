package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pduzinki/fpl-price-checker/internal/di"
	"github.com/pduzinki/fpl-price-checker/internal/domain"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

func main() {
	gs := di.NewGetServiceLambdas()

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Info().Msg(fmt.Sprintf("api gateway request: %v", req))

		var report domain.PriceChangeReport
		var err error

		if date, prs := req.PathParameters["date"]; prs {
			report, err = gs.GetReportByDate(ctx, date)
		} else {
			report, err = gs.GetLatestReport(ctx)
		}

		if err != nil {
			log.Error().Msg(fmt.Sprintf("failed to get report: %v", err))

			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"message": "404 not found"}`,
				StatusCode: http.StatusNotFound,
			}, nil
		}

		body, err := json.Marshal(report)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("failed to marshal report: %v", err))

			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"message": "500 internal server error"}`,
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(body),
			StatusCode: http.StatusOK,
		}, nil
	})
}
