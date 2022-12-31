package fetch

import (
	"context"
	"fmt"
	"time"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"
)

const dateFormat = "2006-01-02"

type PlayersGetter interface {
	GetPlayers() ([]wrapper.Player, error)
}

type StorageAdder interface {
	Add(ctx context.Context, date string, players map[int]domain.Player) error
}

type FetchService struct {
	pg PlayersGetter // TODO find a better name for this interface
	sa StorageAdder  // TODO find a better name for this interface
}

func NewFetchService(pg PlayersGetter, sa StorageAdder) FetchService {
	return FetchService{
		pg: pg,
		sa: sa,
	}
}

func (fs *FetchService) Fetch() error {
	players, err := fs.pg.GetPlayers()
	if err != nil {
		return fmt.Errorf("FetchService.Fetch, failed to get data from api: %w", err)
	}

	playersMap := make(map[int]domain.Player)

	for _, p := range players {
		playersMap[p.ID] = toDomainPlayer(&p)
	}

	todaysDate := time.Now().Format(dateFormat)

	err = fs.sa.Add(context.TODO(), todaysDate, playersMap)
	if err != nil {
		return fmt.Errorf("FetchService.Fetch, failed to save data in storage: %w", err)
	}

	return nil
}

func toDomainPlayer(wp *wrapper.Player) domain.Player {
	return domain.Player{
		ID:         wp.ID,
		Name:       wp.WebName,
		Price:      wp.Price,
		SelectedBy: wp.SelectedBy,
	}
}
