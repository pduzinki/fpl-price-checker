package fs

import (
	"context"
	"fmt"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

// TODO add later

var testPlayers = map[int]domain.Player{
	1: {ID: 1, Name: "Kane", Price: 120, SelectedBy: "12%"},
	2: {ID: 2, Name: "Salah", Price: 130, SelectedBy: "22%"},
	3: {ID: 3, Name: "Haaland", Price: 125, SelectedBy: "80%"},
}

func TestFoo(t *testing.T) {
	ctx := context.Background()

	repo, _ := NewDailyPlayersDataRepository("./data")

	repo.Add(ctx, "2022-01-01", testPlayers)

	players, err := repo.GetByDate(ctx, "2022-01-01")

	fmt.Println(err)

	fmt.Println(players)
}
