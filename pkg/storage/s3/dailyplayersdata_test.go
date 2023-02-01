package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/config"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

func Disabled_TestBar(t *testing.T) {
	// TODO add proper tests

	date := "12-12-12"
	players := domain.DailyPlayersData{
		1: domain.Player{
			ID:         1,
			Name:       "Kane",
			Price:      123,
			SelectedBy: "55.5",
		},
		2: domain.Player{
			ID:         2,
			Name:       "Salah",
			Price:      130,
			SelectedBy: "23.1",
		},
	}

	ctx := context.Background()
	cfg := config.NewConfig()

	pr, err := NewDailyPlayersDataRepository(cfg.AWS, "players")
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	err = pr.Add(ctx, date, players)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	got, err := pr.GetByDate(ctx, date)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	fmt.Println(got)
}
