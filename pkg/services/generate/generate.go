package generate

import (
	"context"
	"fmt"
	"time"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

type DailyPlayersDataGetter interface {
	GetByDate(ctx context.Context, date string) (domain.DailyPlayersData, error)
}

type PriceChangeReportAdder interface {
	Add(ctx context.Context, date string, report domain.PriceChangeReport) error
}

type GenerateService struct {
	dg DailyPlayersDataGetter
	ra PriceChangeReportAdder
}

func NewGenerateService(sg DailyPlayersDataGetter, ra PriceChangeReportAdder) *GenerateService {
	return &GenerateService{
		dg: sg,
		ra: ra,
	}
}

func (gs *GenerateService) GeneratePriceReport(ctx context.Context) error {
	todaysDate := time.Now().Format(domain.DateFormat)
	yesterdaysDate := time.Now().Add(-24 * time.Hour).Format(domain.DateFormat)

	yesterdayPlayers, err := gs.dg.GetByDate(ctx, yesterdaysDate)
	if err != nil {
		return fmt.Errorf("generate.GenerateService.GeneratePriceReport failed to get yesterday's players data: %w", err)
	}

	todayPlayers, err := gs.dg.GetByDate(ctx, todaysDate)
	if err != nil {
		return fmt.Errorf("generate.GenerateService.GeneratePriceReport failed to get today's players data: %w", err)
	}

	report := domain.PriceChangeReport{
		Date:    todaysDate,
		Records: make([]domain.Record, 0),
	}

	priceChangedPlayers := make([]domain.Record, 0)
	newPlayers := make([]domain.Record, 0)

	for tk, tv := range todayPlayers {
		yv, prs := yesterdayPlayers[tk]
		if !prs {
			newPlayers = append(newPlayers, domain.Record{
				Name:        tv.Name,
				OldPrice:    "-",
				NewPrice:    fmt.Sprintf("%.1f", float64(tv.Price)/10.),
				Description: "new",
			})
			continue
		}

		if yv.Price != tv.Price {
			record := domain.Record{
				Name:        tv.Name,
				OldPrice:    fmt.Sprintf("%.1f", float64(yv.Price)/10.),
				NewPrice:    fmt.Sprintf("%.1f", float64(tv.Price)/10.),
				Description: addDescription(yv.Price, tv.Price),
			}

			priceChangedPlayers = append(priceChangedPlayers, record)
		}
	}

	report.Records = append(priceChangedPlayers, newPlayers...)

	err = gs.ra.Add(ctx, todaysDate, report)
	if err != nil {
		return fmt.Errorf("generate.GenerateService.GeneratePriceReport failed to save report: %w", err)
	}

	return nil
}

func addDescription(oldPrice, newPrice int) string {
	if oldPrice > newPrice {
		return "drop"
	}
	return "rise"
}
