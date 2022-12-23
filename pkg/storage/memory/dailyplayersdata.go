package memory

import (
	"context"
	"sync"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"
)

// TODO add type alias for map[int]domain.Player

type DailyPlayersDataRepository struct {
	dailyPlayers map[string](map[int]domain.Player)
	sync.RWMutex
}

func NewDailyPlayersDataRepository() DailyPlayersDataRepository {
	return DailyPlayersDataRepository{
		dailyPlayers: make(map[string]map[int]domain.Player),
	}
}

// TODO add input validations (e.g. if date is a proper date)

func (dr *DailyPlayersDataRepository) Add(_ context.Context, date string, players map[int]domain.Player) error {
	dr.Lock()
	defer dr.Unlock()

	if _, ok := dr.dailyPlayers[date]; ok {
		return storage.ErrDataAlreadyExists
	}
	dr.dailyPlayers[date] = players

	return nil
}

func (dr *DailyPlayersDataRepository) Update(_ context.Context, date string, players map[int]domain.Player) error {
	dr.Lock()
	defer dr.Unlock()

	if _, ok := dr.dailyPlayers[date]; ok {
		dr.dailyPlayers[date] = players
	}

	return storage.ErrDataNotFound
}

func (dr *DailyPlayersDataRepository) GetByDate(_ context.Context, date string) (map[int]domain.Player, error) {
	dr.RLock()
	defer dr.RUnlock()

	if dailyPlayers, ok := dr.dailyPlayers[date]; ok {
		return dailyPlayers, nil
	}

	return nil, storage.ErrDataNotFound
}
