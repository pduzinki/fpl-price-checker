package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/config"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/localstack"
	"github.com/stretchr/testify/suite"
)

type DailyPlayersDataRepositoryTestSuite struct {
	suite.Suite
	c    *gnomock.Container
	repo *DailyPlayersDataRepository
}

func (suite *DailyPlayersDataRepositoryTestSuite) SetupSuite() {
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

	repo, err := NewDailyPlayersDataRepository(cfg, "players")
	suite.NoError(err)

	suite.repo = repo
}

func (suite *DailyPlayersDataRepositoryTestSuite) TearDownSuite() {
	suite.NoError(gnomock.Stop(suite.c))
}

func TestDailyPlayersDataRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DailyPlayersDataRepositoryTestSuite))
}

func (suite *DailyPlayersDataRepositoryTestSuite) TestDailyPlayersDataAddAndGetByDate() {
	ctx := context.Background()

	date := "1999-06-24"
	players := domain.DailyPlayersData{
		1: domain.Player{
			ID:         1,
			Name:       "Haaland",
			Price:      132,
			SelectedBy: "84.3",
		},
	}

	err := suite.repo.Add(ctx, date, players)
	suite.NoError(err)

	gotPlayers, err := suite.repo.GetByDate(ctx, date)
	suite.NoError(err)
	suite.EqualValues(players, gotPlayers)
}

func (suite *DailyPlayersDataRepositoryTestSuite) TestDailyPlayersDataAddDuplicate() {
	ctx := context.Background()

	date := "1999-06-25"
	players := make(domain.DailyPlayersData)

	err := suite.repo.Add(ctx, date, players)
	suite.NoError(err)

	err = suite.repo.Add(ctx, date, players)
	suite.ErrorIs(err, storage.ErrDataAlreadyExists)
}

func (suite *DailyPlayersDataRepositoryTestSuite) TestDailyPlayersDataGetNonExistentEntry() {
	ctx := context.Background()

	date := "1999-06-26"
	players, err := suite.repo.GetByDate(ctx, date)
	suite.ErrorIs(err, storage.ErrDataNotFound)
	suite.Nil(players)
}
