package generate

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pduzinki/fpl-price-checker/internal/domain"

	"github.com/stretchr/testify/assert"
)

var (
	errDataGetterFailure  = fmt.Errorf("data getter failure")
	errReportAdderFailure = fmt.Errorf("report adder failure")
	errTeamsGetterFailure = fmt.Errorf("teams getter failure")

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

	teamsGetterOk TeamsGetter = &TeamsGetterMock{
		GetAllFunc: func() (map[int]domain.Team, error) {
			return map[int]domain.Team{}, nil
		},
	}

	dataGetterFailing DailyPlayersDataGetter = &DailyPlayersDataGetterMock{
		GetByDateFunc: func(ctx context.Context, date string) (domain.DailyPlayersData, error) {
			return domain.DailyPlayersData{}, errDataGetterFailure
		},
	}

	reportAdderFailing PriceChangeReportAdder = &PriceChangeReportAdderMock{
		AddFunc: func(ctx context.Context, date string, report domain.PriceChangeReport) error {
			return errReportAdderFailure
		},
	}

	teamsGetterFailing TeamsGetter = &TeamsGetterMock{
		GetAllFunc: func() (map[int]domain.Team, error) {
			return nil, errTeamsGetterFailure
		},
	}
)

func TestGenerate(t *testing.T) {
	testcases := []struct {
		name    string
		dg      DailyPlayersDataGetter
		ra      PriceChangeReportAdder
		tg      TeamsGetter
		wantErr error
	}{
		{
			name:    "sunny scenario",
			dg:      dataGetterOk,
			ra:      reportAdderOk,
			tg:      teamsGetterOk,
			wantErr: nil,
		},
		{
			name:    "DailyPlayersDataGetter failure",
			dg:      dataGetterFailing,
			ra:      reportAdderOk,
			tg:      teamsGetterOk,
			wantErr: errDataGetterFailure,
		},
		{
			name:    "PriceChangeReportAdder failure",
			dg:      dataGetterOk,
			ra:      reportAdderFailing,
			tg:      teamsGetterOk,
			wantErr: errReportAdderFailure,
		},
		{
			name:    "TeamsGetter failure",
			dg:      dataGetterOk,
			ra:      reportAdderFailing,
			tg:      teamsGetterFailing,
			wantErr: errTeamsGetterFailure,
		},
	}

	for _, test := range testcases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			g := NewGenerateService(test.dg, test.ra, test.tg)

			err := g.GeneratePriceReport(ctx)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("want: %v, got: %v", test.wantErr, err)
			}
		})
	}
}

var (
	playersData1 = domain.DailyPlayersData{
		1: domain.Player{
			ID:         1,
			TeamID:     1,
			Name:       "Haaland",
			Price:      132,
			SelectedBy: "84.3",
		},
		2: domain.Player{
			ID:         2,
			TeamID:     2,
			Name:       "Salah",
			Price:      119,
			SelectedBy: "20.1",
		},
		3: domain.Player{
			ID:         3,
			TeamID:     3,
			Name:       "Kane",
			Price:      105,
			SelectedBy: "5.1",
		},
	}
	playersData2 = domain.DailyPlayersData{
		1: domain.Player{
			ID:         1,
			TeamID:     1,
			Name:       "Haaland",
			Price:      133,
			SelectedBy: "84.3",
		},
		2: domain.Player{
			ID:         2,
			TeamID:     2,
			Name:       "Salah",
			Price:      119,
			SelectedBy: "20.1",
		},
		3: domain.Player{
			ID:         3,
			TeamID:     3,
			Name:       "Kane",
			Price:      104,
			SelectedBy: "5.1",
		},
	}
	playersData3 = domain.DailyPlayersData{
		1: domain.Player{
			ID:         1,
			TeamID:     1,
			Name:       "Haaland",
			Price:      133,
			SelectedBy: "84.3",
		},
		2: domain.Player{
			ID:         2,
			TeamID:     2,
			Name:       "Salah",
			Price:      119,
			SelectedBy: "20.1",
		},
		3: domain.Player{
			ID:         3,
			TeamID:     3,
			Name:       "Kane",
			Price:      104,
			SelectedBy: "5.1",
		},
		4: domain.Player{
			ID:         4,
			TeamID:     4,
			Name:       "Mitoma",
			Price:      50,
			SelectedBy: "0.0",
		},
	}

	records1 = []domain.Record{
		{
			Name:        "Haaland",
			Team:        "MCI",
			OldPrice:    "13.2",
			NewPrice:    "13.3",
			Description: "rise",
		},
		{
			Name:        "Kane",
			Team:        "TOT",
			OldPrice:    "10.5",
			NewPrice:    "10.4",
			Description: "drop",
		},
	}

	records2 = []domain.Record{
		{
			Name:        "Haaland",
			Team:        "MCI",
			OldPrice:    "-",
			NewPrice:    "13.2",
			Description: "new",
		},
		{
			Name:        "Salah",
			Team:        "LIV",
			OldPrice:    "-",
			NewPrice:    "11.9",
			Description: "new",
		},
		{
			Name:        "Kane",
			Team:        "TOT",
			OldPrice:    "-",
			NewPrice:    "10.5",
			Description: "new",
		},
	}

	records3 = []domain.Record{
		{
			Name:        "Haaland",
			Team:        "MCI",
			OldPrice:    "13.2",
			NewPrice:    "13.3",
			Description: "rise",
		},
		{
			Name:        "Kane",
			Team:        "TOT",
			OldPrice:    "10.5",
			NewPrice:    "10.4",
			Description: "drop",
		},
		{
			Name:        "Mitoma",
			Team:        "BHA",
			OldPrice:    "-",
			NewPrice:    "5.0",
			Description: "new",
		},
	}

	teams1 = domain.Teams{
		1: {
			ID:        1,
			Shortname: "MCI",
		},
		2: {
			ID:        2,
			Shortname: "LIV",
		},
		3: {
			ID:        3,
			Shortname: "TOT",
		},
		4: {
			ID:        4,
			Shortname: "BHA",
		},
	}
)

func TestGenerateRecords(t *testing.T) {
	testcases := []struct {
		name             string
		yesterdayPlayers domain.DailyPlayersData
		todayPlayers     domain.DailyPlayersData
		teams            domain.Teams
		want             []domain.Record
	}{
		{
			name:             "price changes for some players",
			yesterdayPlayers: playersData1,
			todayPlayers:     playersData2,
			teams:            teams1,
			want:             records1,
		},
		{
			name:             "new players added",
			yesterdayPlayers: domain.DailyPlayersData{},
			todayPlayers:     playersData1,
			teams:            teams1,
			want:             records2,
		},
		{
			name:             "price changes and new players",
			yesterdayPlayers: playersData1,
			todayPlayers:     playersData3,
			teams:            teams1,
			want:             records3,
		},
		{
			name:             "no price changes, no new players",
			yesterdayPlayers: playersData1,
			todayPlayers:     playersData1,
			teams:            teams1,
			want:             []domain.Record{},
		},
	}

	for _, test := range testcases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			got := generateRecords(test.yesterdayPlayers, test.todayPlayers, test.teams)
			assert.ElementsMatch(t, got, test.want)
		})
	}
}
