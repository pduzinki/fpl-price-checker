package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/config"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

func TestFoo(t *testing.T) {
	// TODO replace with proper tests

	report := domain.PriceChangeReport{
		Date: "12-12-12",
		Records: []domain.Record{
			{
				Name:        "Kane",
				OldPrice:    "12.2",
				NewPrice:    "12.3",
				Description: "rise",
			},
		},
	}

	ctx := context.Background()
	cfg, _ := config.NewConfig()

	pr, err := NewPriceReportRepository(cfg.AWS, "reports")
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	err = pr.Add(ctx, report.Date, report)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	got, err := pr.GetByDate(ctx, report.Date)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	fmt.Println(got)
}
