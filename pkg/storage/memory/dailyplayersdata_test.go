package memory

import (
	"context"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"
)

var testPlayers = map[int]domain.Player{
	1: {ID: 1, Name: "Kane", Price: 120, SelectedBy: "12%"},
	2: {ID: 2, Name: "Salah", Price: 130, SelectedBy: "22%"},
	3: {ID: 3, Name: "Haaland", Price: 125, SelectedBy: "80%"},
}

func TestDailyPlayersDataRepositoryAdd(t *testing.T) {
	testcases := []struct {
		name    string
		date    string
		players map[int]domain.Player
		want    error
	}{
		{
			name:    "sunny scenario",
			date:    "2022-01-01",
			players: testPlayers,
			want:    nil,
		},
		{
			name:    "data from that date already exists",
			date:    "2022-03-14",
			players: testPlayers,
			want:    storage.ErrDataAlreadyExists,
		},
		// TODO add case "date format wrong"
	}

	dailyPlayersDataRepo := DailyPlayersDataRepository{
		dailyPlayers: map[string]map[int]domain.Player{
			"2022-03-14": testPlayers,
		},
	}

	for _, test := range testcases {
		test := test

		got := dailyPlayersDataRepo.Add(context.Background(), test.date, test.players)
		if got != test.want {
			t.Errorf("want: %v, got: %v", test.want, got)
		}

		if got == nil {
			continue
		}

		// TODO compare values from test.data with what was saved in map
	}
}

func TestDailyPlayersDataRepositoryUpdate(t *testing.T) {
	// TODO add test

	testcases := []struct {
		name string
	}{
		{},
	}

	for _, test := range testcases {
		test := test

		_ = test
	}
}

func TestDailyPlayersDataRepositoryGetByDate(t *testing.T) {
	// TODO add test

	testcases := []struct {
		name string
	}{
		{},
	}

	for _, test := range testcases {
		test := test

		_ = test
	}
}
