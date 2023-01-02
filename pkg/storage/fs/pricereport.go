package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"
)

type PriceReportRepository struct {
	folderPath string
	sync.RWMutex
}

func NewPriceReportRepository(folderPath string) (*PriceReportRepository, error) {
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create price report data repository: %w", err)
	}

	return &PriceReportRepository{
		folderPath: folderPath,
	}, nil
}

func (pr *PriceReportRepository) Add(_ context.Context, date string, report domain.PriceChangeReport) error {
	pr.Lock()
	defer pr.Unlock()

	if err := domain.ParseDate(date); err != nil {
		return fmt.Errorf("fs.PriceReportRepository.Add failed to parse date: %w", err)
	}

	filename := filepath.Join(pr.folderPath, date)

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("fs.PriceReportRepository.Add failed: %w", storage.ErrDataAlreadyExists)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("fs.PriceReportRepository.Add failed to create file: %w", err)
	}
	defer f.Close()

	jsonReport, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("fs.PriceReportRepository.Add failed to marshal data: %w", err)
	}

	if _, err := f.Write(jsonReport); err != nil {
		return fmt.Errorf("fs.PriceReportRepository.Add failed to write data into file: %w", err)
	}

	return nil
}

func (pr *PriceReportRepository) GetByDate(_ context.Context, date string) (*domain.PriceChangeReport, error) {
	pr.Lock()
	defer pr.Unlock()

	// TODO

	return nil, fmt.Errorf("unimplemented")
}
