package fs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"
)

type DailyPlayersDataRepository struct {
	folderPath string
	sync.RWMutex
}

func NewDailyPlayersDataRepository(folderPath string) (*DailyPlayersDataRepository, error) {
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create daily players data repository: %w", err)
	}

	return &DailyPlayersDataRepository{
		folderPath: folderPath,
	}, nil
}

func (dr *DailyPlayersDataRepository) Add(_ context.Context, date string, players map[int]domain.Player) error {
	dr.Lock()
	defer dr.Unlock()

	// TODO check if date is proper date

	filename := filepath.Join(dr.folderPath, date)

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("fs.Add failed: %w", storage.ErrDataAlreadyExists)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("fs.Add, failed to create file: %w", err)
	}
	defer f.Close()

	jsonPlayers, err := json.Marshal(players)
	if err != nil {
		return fmt.Errorf("fs.Add, failed to marshal data: %w", err)
	}

	if _, err := f.Write(jsonPlayers); err != nil {
		return fmt.Errorf("fs.Add, failed to write data into file: %w", err)
	}

	return nil
}

func (dr *DailyPlayersDataRepository) Update(_ context.Context, date string, players map[int]domain.Player) error {
	dr.Lock()
	defer dr.Unlock()

	return errors.New("unimplemented")
}

func (dr *DailyPlayersDataRepository) GetByDate(_ context.Context, date string) (map[int]domain.Player, error) {
	dr.RLock()
	defer dr.RUnlock()

	filename := filepath.Join(dr.folderPath, date)

	if _, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("fs.GetByDate, failed to fetch file info: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("fs.GetByDate, failed to read file: %w", err)
	}

	players := make(map[int]domain.Player, 0)

	if err := json.Unmarshal(data, &players); err != nil {
		return nil, fmt.Errorf("fs.GetByDate, failed to unmarshal data: %w", err)
	}

	return players, nil
}
