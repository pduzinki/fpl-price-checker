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

type Player struct {
	ID         int
	Name       string
	Price      int
	SelectedBy string
}

type DailyPlayersData map[int]Player

type DailyPlayersDataRepository struct {
	folderPath string
	sync.RWMutex
}

func NewDailyPlayersDataRepository(folderPath string) (*DailyPlayersDataRepository, error) {
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("fs.NewDailyPlayersDataRepository failed: %w", err)
	}

	return &DailyPlayersDataRepository{
		folderPath: folderPath,
	}, nil
}

func (dr *DailyPlayersDataRepository) Add(_ context.Context, date string, players domain.DailyPlayersData) error {
	dr.Lock()
	defer dr.Unlock()

	if err := domain.ParseDate(date); err != nil {
		return fmt.Errorf("fs.DailyPlayersDataRepository.Add failed to parse date: %w", err)
	}

	filename := filepath.Join(dr.folderPath, date)

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("fs.DailyPlayersDataRepository.Add failed: %w", storage.ErrDataAlreadyExists)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("fs.DailyPlayersDataRepository.Add, failed to create file: %w", err)
	}
	defer f.Close()

	fsPlayers := toFsDailyPlayersData(players)

	jsonPlayers, err := json.Marshal(fsPlayers)
	if err != nil {
		return fmt.Errorf("fs.DailyPlayersDataRepository.Add, failed to marshal data: %w", err)
	}

	if _, err := f.Write(jsonPlayers); err != nil {
		return fmt.Errorf("fs.DailyPlayersDataRepository.Add, failed to write data into file: %w", err)
	}

	return nil
}

func (dr *DailyPlayersDataRepository) Update(_ context.Context, date string, players domain.DailyPlayersData) error {
	dr.Lock()
	defer dr.Unlock()

	return errors.New("unimplemented")
}

func (dr *DailyPlayersDataRepository) GetByDate(_ context.Context, date string) (domain.DailyPlayersData, error) {
	dr.RLock()
	defer dr.RUnlock()

	filename := filepath.Join(dr.folderPath, date)

	if _, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("fs.DailyPlayersDataRepository.GetByDate, failed to fetch file info: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("fs.DailyPlayersDataRepository.GetByDate, failed to read file: %w", err)
	}

	fsPlayers := make(DailyPlayersData, 0)

	if err := json.Unmarshal(data, &fsPlayers); err != nil {
		return nil, fmt.Errorf("fs.DailyPlayersDataRepository.GetByDate, failed to unmarshal data: %w", err)
	}

	return toDomainDailyPlayersData(fsPlayers), nil
}

func toDomainDailyPlayersData(data DailyPlayersData) domain.DailyPlayersData {
	domainDailyPlayersData := make(domain.DailyPlayersData)

	for k, v := range data {
		domainDailyPlayersData[k] = domain.Player(v)
	}

	return domainDailyPlayersData
}

func toFsDailyPlayersData(data domain.DailyPlayersData) DailyPlayersData {
	fsDailyPlayersData := make(DailyPlayersData)

	for k, v := range data {
		fsDailyPlayersData[k] = Player(v)
	}

	return fsDailyPlayersData
}
