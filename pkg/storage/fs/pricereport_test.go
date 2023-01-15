package fs

import (
	"context"
	"os"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"

	"github.com/stretchr/testify/suite"
)

type PriceReportRepositoryTestSuite struct {
	suite.Suite
	folderPath string
	repo       *PriceReportRepository
}

func (suite *PriceReportRepositoryTestSuite) SetupSuite() {
	suite.folderPath = "./tmp_price_report_repository_test_data/"

	err := os.MkdirAll(suite.folderPath, 0755)
	suite.NoError(err)

	repo, err := NewPriceReportRepository(suite.folderPath)
	suite.NoError(err)

	suite.repo = repo
}

func (suite *PriceReportRepositoryTestSuite) TearDownSuite() {
	err := os.RemoveAll(suite.folderPath)
	suite.NoError(err)
}

func TestPriceReportTestSuite(t *testing.T) {
	suite.Run(t, new(PriceReportRepositoryTestSuite))
}

func (suite *PriceReportRepositoryTestSuite) TestPriceReportAddAndGetByDate() {
	ctx := context.Background()

	date := "1999-06-24"
	report := domain.PriceChangeReport{
		Date: date,
		Records: []domain.Record{
			{
				Name:        "Kane",
				OldPrice:    "12.2",
				NewPrice:    "12.3",
				Description: "rise",
			},
		},
	}

	err := suite.repo.Add(ctx, date, report)
	suite.NoError(err)

	gotReport, err := suite.repo.GetByDate(ctx, date)
	suite.NoError(err)
	suite.EqualValues(report, gotReport)
}

func (suite *PriceReportRepositoryTestSuite) TestPriceReportAddDuplicate() {
	ctx := context.Background()

	date := "1999-06-25"
	report := domain.PriceChangeReport{
		Date: date,
		Records: []domain.Record{
			{
				Name:        "Salah",
				OldPrice:    "12.2",
				NewPrice:    "12.1",
				Description: "drop",
			},
		},
	}

	err := suite.repo.Add(ctx, date, report)
	suite.NoError(err)

	err = suite.repo.Add(ctx, date, report)
	suite.ErrorIs(err, storage.ErrDataAlreadyExists)
}

func (suite *PriceReportRepositoryTestSuite) TestPriceReportGetNonExistentEntry() {
	ctx := context.Background()

	date := "1999-06-26"
	report, err := suite.repo.GetByDate(ctx, date)
	suite.ErrorIs(err, storage.ErrDataNotFound)
	suite.EqualValues(report, domain.PriceChangeReport{})
}

func (suite *PriceReportRepositoryTestSuite) TestPriceReportAddWithIncorrectlyFormattedDate() {
	ctx := context.Background()

	date := "not-even-a-date"
	report := domain.PriceChangeReport{
		Date: date,
		Records: []domain.Record{
			{
				Name:        "Kane",
				OldPrice:    "12.2",
				NewPrice:    "12.3",
				Description: "rise",
			},
		},
	}

	err := suite.repo.Add(ctx, date, report)
	suite.Error(err)
}
