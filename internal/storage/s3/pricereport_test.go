package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/pduzinki/fpl-price-checker/internal/config"
	"github.com/pduzinki/fpl-price-checker/internal/domain"
	"github.com/pduzinki/fpl-price-checker/internal/storage"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/localstack"
	"github.com/stretchr/testify/suite"
)

type PriceReportRepositoryTestSuite struct {
	suite.Suite
	c    *gnomock.Container
	repo *PriceReportRepository
}

func (suite *PriceReportRepositoryTestSuite) SetupSuite() {
	p := localstack.Preset(
		localstack.WithServices(localstack.S3),
		localstack.WithS3Files("testdata"),
		localstack.WithVersion("0.12.0"),
	)

	c, err := gnomock.Start(p)
	suite.NoError(err)

	suite.c = c

	cfg := config.AWSConfig{
		Region:   "eu-west-2",
		ID:       "test",
		Secret:   "test",
		Endpoint: fmt.Sprintf("http://%s/", c.Address(localstack.APIPort)),
		Bucket:   "fpc-bucket",
	}

	repo, err := NewPriceReportRepository(cfg, "reports")
	suite.NoError(err)

	suite.repo = repo
}

func (suite *PriceReportRepositoryTestSuite) TearDownSuite() {
	suite.NoError(gnomock.Stop(suite.c))
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
