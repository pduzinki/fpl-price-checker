package fs

import (
	"context"
	"os"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"
	"github.com/stretchr/testify/suite"
)

type DailyPlayersDataRepositoryTestSuite struct {
	suite.Suite
	folderPath string
	repo       *DailyPlayersDataRepository
}

func (suite *DailyPlayersDataRepositoryTestSuite) SetupSuite() {
	suite.folderPath = "./tmp_daily_players_data_test_data/"

	err := os.MkdirAll(suite.folderPath, 0755)
	suite.NoError(err)

	repo, err := NewDailyPlayersDataRepository(suite.folderPath)
	suite.NoError(err)

	suite.repo = repo
}

func (suite *DailyPlayersDataRepositoryTestSuite) TearDownSuite() {
	err := os.RemoveAll(suite.folderPath)
	suite.NoError(err)
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

func (suite *DailyPlayersDataRepositoryTestSuite) TestDailyPlayersDataAddWithIncorrectlyFormattedDate() {
	ctx := context.Background()

	date := "not-a-date"
	players := domain.DailyPlayersData{
		1: domain.Player{
			ID:         1,
			Name:       "Haaland",
			Price:      132,
			SelectedBy: "84.3",
		},
	}

	err := suite.repo.Add(ctx, date, players)
	suite.Error(err)
}
