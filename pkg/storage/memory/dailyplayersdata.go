package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"
)

type DailyPlayersDataRepository struct {
	dailyPlayers map[string]([]domain.Player) // TODO change []domain.Player to map[ID int]domain.Player
	sync.RWMutex
}

func NewDailyPlayersDataRepository() DailyPlayersDataRepository {
	return DailyPlayersDataRepository{
		dailyPlayers: make(map[string][]domain.Player),
	}
}

// TODO add input validations (e.g. if date is a proper date)

func (dr *DailyPlayersDataRepository) Add(_ context.Context, date string, players []domain.Player) error {
	dr.Lock()
	defer dr.Unlock()

	if _, ok := dr.dailyPlayers[date]; ok {
		return storage.ErrDataAlreadyExists
	}

	return nil
}

func (dr *DailyPlayersDataRepository) Update(_ context.Context, date string, players []domain.Player) error {
	dr.Lock()
	defer dr.Unlock()

	if _, ok := dr.dailyPlayers[date]; ok {
		dr.dailyPlayers[date] = players
	}

	return errors.New("data not found")
}

func (dr *DailyPlayersDataRepository) GetByDate(_ context.Context, date string) ([]domain.Player, error) {
	dr.RLock()
	defer dr.RUnlock()

	if dailyPlayers, ok := dr.dailyPlayers[date]; ok {
		return dailyPlayers, nil
	}

	return nil, errors.New("data not found")
}
