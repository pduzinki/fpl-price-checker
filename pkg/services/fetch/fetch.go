package fetch

import (
	"context"
	"fmt"
	"time"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"
)

type PlayersGetter interface {
	GetPlayers() ([]wrapper.Player, error)
}

type DailyPlayersDataAdder interface {
	Add(ctx context.Context, date string, players domain.DailyPlayersData) error
}

type FetchService struct {
	pg PlayersGetter
	da DailyPlayersDataAdder
}

func NewFetchService(pg PlayersGetter, sa DailyPlayersDataAdder) *FetchService {
	return &FetchService{
		pg: pg,
		da: sa,
	}
}

func (fs *FetchService) Fetch(ctx context.Context) error {
	players, err := fs.pg.GetPlayers()
	if err != nil {
		return fmt.Errorf("FetchService.Fetch, failed to get data from api: %w", err)
	}

	playersMap := make(domain.DailyPlayersData)

	for _, p := range players {
		p, err := toDomainPlayer(&p)
		if err != nil {
			return fmt.Errorf("fetch.Fetch failed to convert player value: %w", err)
		}
		playersMap[p.ID] = p
	}

	todaysDate := time.Now().Format(domain.DateFormat)

	err = fs.da.Add(ctx, todaysDate, playersMap)
	if err != nil {
		return fmt.Errorf("FetchService.Fetch, failed to save data in storage: %w", err)
	}

	return nil
}

func toDomainPlayer(wp *wrapper.Player) (domain.Player, error) {
	// selectedBy, err := strconv.ParseFloat(wp.SelectedBy, 32)
	// if err != nil {
	// 	return domain.Player{}, fmt.Errorf("fetch.toDomainPlayer failed to parse float value: %w", err)
	// }

	return domain.Player{
		ID:         wp.ID,
		Name:       wp.WebName,
		Price:      wp.Price,
		SelectedBy: wp.SelectedBy,
	}, nil
}
