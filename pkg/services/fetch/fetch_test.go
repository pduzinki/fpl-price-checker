package fetch

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"
)

var (
	errPlayerGetterFailure = fmt.Errorf("player getter failure")
	errDataAdderFailure    = fmt.Errorf("data adder failure")

	PlayersGetterOk = PlayersGetterMock{
		GetPlayersFunc: func() ([]wrapper.Player, error) {
			return []wrapper.Player{
				{
					ID:         1,
					WebName:    "Kane",
					Price:      123,
					SelectedBy: "33.5",
				},
				{
					ID:         2,
					WebName:    "Salah",
					Price:      130,
					SelectedBy: "55.1",
				},
			}, nil
		},
	}

	PlayersGetterFailing = PlayersGetterMock{
		GetPlayersFunc: func() ([]wrapper.Player, error) {
			return nil, errPlayerGetterFailure
		},
	}

	DailyPlayersDataAdderOk = DailyPlayersDataAdderMock{
		AddFunc: func(ctx context.Context, date string, players domain.DailyPlayersData) error {
			return nil
		},
	}
	DailyPlayersDataAdderFailing = DailyPlayersDataAdderMock{
		AddFunc: func(ctx context.Context, date string, players domain.DailyPlayersData) error {
			return errDataAdderFailure
		},
	}
)

func TestFetch(t *testing.T) {
	testcases := []struct {
		name    string
		pg      PlayersGetter
		da      DailyPlayersDataAdder
		wantErr error
	}{
		{
			name:    "sunny scenario",
			pg:      &PlayersGetterOk,
			da:      &DailyPlayersDataAdderOk,
			wantErr: nil,
		},
		{
			name:    "PlayersGetter failure",
			pg:      &PlayersGetterFailing,
			da:      &DailyPlayersDataAdderOk,
			wantErr: errPlayerGetterFailure,
		},
		{
			name:    "DailyPlayersDataAdder failure",
			pg:      &PlayersGetterOk,
			da:      &DailyPlayersDataAdderFailing,
			wantErr: errDataAdderFailure,
		},
	}

	for _, test := range testcases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			f := NewFetchService(test.pg, test.da)

			err := f.Fetch(ctx)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("want: %v, got: %v", test.wantErr, err)
			}
		})
	}
}
