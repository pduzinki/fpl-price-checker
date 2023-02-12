package get

import (
	"context"
	"fmt"
	"time"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"

	"go.uber.org/multierr"
)

//go:generate moq -out get_moq_test.go . PriceChangeReportGetter

type PriceChangeReportGetter interface {
	GetByDate(ctx context.Context, date string) (domain.PriceChangeReport, error)
}

type GetService struct {
	rg PriceChangeReportGetter
}

func NewGetService(rg PriceChangeReportGetter) *GetService {
	return &GetService{
		rg: rg,
	}
}

func (gs *GetService) GetLatestReport(ctx context.Context) (domain.PriceChangeReport, error) {
	todaysDate := time.Now().Format(domain.DateFormat)
	yesterdaysDate := time.Now().Add(-24 * time.Hour).Format(domain.DateFormat)

	var errors error

	// NOTE: FPL price changes usually occur around 1:30am GMT, so there's a time gap,
	// where there's no "today's" report just yet. in that case, latest will be a report from the day before.

	report, err := gs.rg.GetByDate(ctx, todaysDate)
	if err == nil {
		return report, nil
	}
	errors = multierr.Append(errors, err)

	report, err = gs.rg.GetByDate(ctx, yesterdaysDate)
	if err == nil {
		return report, nil
	}
	errors = multierr.Append(errors, err)

	return domain.PriceChangeReport{}, fmt.Errorf("get.GetService.GetLatestReport failed: %w", errors)
}

func (gs *GetService) GetReportByDate(ctx context.Context, date string) (domain.PriceChangeReport, error) {
	if err := domain.ParseDate(date); err != nil {
		return domain.PriceChangeReport{}, err
	}

	report, err := gs.rg.GetByDate(ctx, date)
	if err != nil {
		return domain.PriceChangeReport{}, fmt.Errorf("get.GetService.GetReportByDate failed: %w", err)
	}

	return report, nil
}
