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

type Record struct {
	Name        string
	OldPrice    string
	NewPrice    string
	Description string
}

type PriceChangeReport struct {
	Date    string
	Records []Record
}

type PriceReportRepository struct {
	folderPath string
	sync.RWMutex
}

func NewPriceReportRepository(folderPath string) (*PriceReportRepository, error) {
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("fs.NewPriceReportRepository failed: %w", err)
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

	fsReport := toFsReport(report)

	jsonReport, err := json.Marshal(fsReport)
	if err != nil {
		return fmt.Errorf("fs.PriceReportRepository.Add failed to marshal data: %w", err)
	}

	if _, err := f.Write(jsonReport); err != nil {
		return fmt.Errorf("fs.PriceReportRepository.Add failed to write data into file: %w", err)
	}

	return nil
}

func (pr *PriceReportRepository) GetByDate(_ context.Context, date string) (domain.PriceChangeReport, error) {
	pr.Lock()
	defer pr.Unlock()

	var report domain.PriceChangeReport

	filename := filepath.Join(pr.folderPath, date)

	if _, err := os.Stat(filename); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = storage.ErrDataNotFound
		}

		return report, fmt.Errorf("fs.PriceReportRepository.GetByDate, failed to fetch file info: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return report, fmt.Errorf("fs.PriceReportRepository.GetByDate, failed to read file: %w", err)
	}

	var fsReport PriceChangeReport

	if err := json.Unmarshal(data, &fsReport); err != nil {
		return report, fmt.Errorf("fs.PriceReportRepository.GetByDate, failed to unmarshal data: %w", err)
	}

	report = toDomainReport(fsReport)

	return report, nil
}

func toFsReport(report domain.PriceChangeReport) PriceChangeReport {
	records := make([]Record, 0, len(report.Records))

	for _, r := range report.Records {
		records = append(records, Record(r))
	}

	return PriceChangeReport{
		Date:    report.Date,
		Records: records,
	}
}

func toDomainReport(report PriceChangeReport) domain.PriceChangeReport {
	records := make([]domain.Record, 0, len(report.Records))

	for _, r := range report.Records {
		records = append(records, domain.Record(r))
	}

	return domain.PriceChangeReport{
		Date:    report.Date,
		Records: records,
	}
}
