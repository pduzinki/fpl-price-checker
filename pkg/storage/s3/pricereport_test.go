package s3

import (
	"context"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/config"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"

	"github.com/elgohr/go-localstack"
	"github.com/stretchr/testify/suite"
)

type PriceReportRepositoryTestSuite struct {
	suite.Suite
	l    *localstack.Instance
	repo *PriceReportRepository
}

func (suite *PriceReportRepositoryTestSuite) SetupSuite() {
	l, err := localstack.NewInstance()
	suite.NoError(err)

	suite.l = l

	err = suite.l.Start()
	suite.NoError(err)

	cfg := config.AWSConfig{
		Region: "eu-west-2",
	}

	repo, err := NewPriceReportRepository(cfg, "reports")
	suite.NoError(err)

	// TODO create bucket

	suite.repo = repo
}

func (suite *PriceReportRepositoryTestSuite) TearDownSuite() {
	err := suite.l.Stop()
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
