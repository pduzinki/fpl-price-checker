package generate

import (
	"context"
	"errors"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

var (
	dataGetterOk DailyPlayersDataGetter = &DailyPlayersDataGetterMock{
		GetByDateFunc: func(ctx context.Context, date string) (domain.DailyPlayersData, error) {
			return domain.DailyPlayersData{
				1: domain.Player{
					ID:         1,
					Name:       "Kane",
					Price:      123,
					SelectedBy: "12.3",
				},
				2: domain.Player{
					ID:         2,
					Name:       "Salah",
					Price:      130,
					SelectedBy: "55.1",
				},
			}, nil
		},
	}

	reportAdderOk PriceChangeReportAdder = &PriceChangeReportAdderMock{
		AddFunc: func(ctx context.Context, date string, report domain.PriceChangeReport) error {
			return nil
		},
	}
)

func TestGenerate(t *testing.T) {
	testcases := []struct {
		name    string
		dg      DailyPlayersDataGetter
		ra      PriceChangeReportAdder
		wantErr error
	}{
		{
			name:    "sunny scenario",
			dg:      dataGetterOk,
			ra:      reportAdderOk,
			wantErr: nil,
		},
		// TODO add more testcases
	}

	for _, test := range testcases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			g := NewGenerateService(test.dg, test.ra)

			err := g.GeneratePriceReport(ctx)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("want: %v, got: %v", test.wantErr, err)
			}
		})
	}
}

func TestGenerateRecords(t *testing.T) {
	// TODO

	testcases := []struct {
		name             string
		yesterdayPlayers domain.DailyPlayersData
		todayPlayers     domain.DailyPlayersData
	}{
		{
			name: "",
		},
	}

	for _, test := range testcases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			got := generateRecords(test.yesterdayPlayers, test.todayPlayers)
			_ = got
		})
	}
}
