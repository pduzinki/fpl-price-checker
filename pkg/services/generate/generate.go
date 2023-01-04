package generate

import (
	"context"
	"fmt"
	"time"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

type StorageGetter interface {
	GetByDate(_ context.Context, date string) (map[int]domain.Player, error)
}

type ReportAdder interface {
	Add(_ context.Context, date string, report domain.PriceChangeReport) error
}

type GenerateService struct {
	sg StorageGetter // TODO find a better name
	ra ReportAdder   // TODO find a better name
}

func NewGenerateService(sg StorageGetter, ra ReportAdder) *GenerateService {
	return &GenerateService{
		sg: sg,
		ra: ra,
	}
}

func (gs *GenerateService) GeneratePriceReport() error {
	todaysDate := time.Now().Format(domain.DateFormat)
	yesterdaysDate := time.Now().Add(-24 * time.Hour).Format(domain.DateFormat)

	yesterdayPlayers, err := gs.sg.GetByDate(context.TODO(), yesterdaysDate)
	if err != nil {
		return fmt.Errorf("err 1: %w", err) // TODO
	}

	todayPlayers, err := gs.sg.GetByDate(context.TODO(), todaysDate)
	if err != nil {
		return fmt.Errorf("err 2: %w", err) // TODO
	}

	report := domain.PriceChangeReport{
		Records: make([]domain.Record, 0),
	}

	for tk, tv := range todayPlayers {
		yv, prs := yesterdayPlayers[tk]
		if !prs {
			// fmt.Println("new:", tv.Name, tv.Price) // TODO no sure if newly added players should be reported
			continue
		}

		if yv.Price != tv.Price {
			record := domain.Record{
				Name:        tv.Name,
				OldPrice:    float64(yv.Price) / 10.,
				NewPrice:    float64(tv.Price) / 10.,
				Description: addDescription(yv.Price, tv.Price),
			}

			fmt.Println(record)

			report.Records = append(report.Records, record)
		}
	}

	err = gs.ra.Add(context.TODO(), todaysDate, report)
	if err != nil {
		return fmt.Errorf("err 33: %w", err) // TODO
	}

	return nil
}

func addDescription(oldPrice, newPrice int) string {
	if oldPrice > newPrice {
		return "drop"
	}
	return "rise"
}
