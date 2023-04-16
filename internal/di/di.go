package di

import (
	"os/user"
	"path/filepath"

	"github.com/pduzinki/fpl-price-checker/internal/config"
	"github.com/pduzinki/fpl-price-checker/internal/domain"
	"github.com/pduzinki/fpl-price-checker/internal/rest"
	"github.com/pduzinki/fpl-price-checker/internal/services/fetch"
	"github.com/pduzinki/fpl-price-checker/internal/services/generate"
	"github.com/pduzinki/fpl-price-checker/internal/services/get"
	"github.com/pduzinki/fpl-price-checker/internal/storage/fs"
	"github.com/pduzinki/fpl-price-checker/internal/storage/memory"
	"github.com/pduzinki/fpl-price-checker/internal/storage/s3"
	"github.com/pduzinki/fpl-price-checker/internal/wrapper"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func Wrapper() *wrapper.Wrapper {
	return wrapper.NewWrapper()
}

func Config() *config.Config {
	cfg := config.NewConfig()

	return cfg
}

func fsDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("di.FsDir failed to get current os user")
	}

	return filepath.Join("/home", user.Username, "fpc", "data")
}

func DailyPlayersDataFsRepository() *fs.DailyPlayersDataRepository {
	dr, err := fs.NewDailyPlayersDataRepository(filepath.Join(fsDir(), "players"))
	if err != nil {
		log.Fatal().Err(err).Msg("di.DailyPlayersDataFsRepository failed")
	}

	return dr
}

func DailyPlayersDataS3Repository() *s3.DailyPlayersDataRepository {
	cfg := Config()

	dr, err := s3.NewDailyPlayersDataRepository(cfg.AWS, "players")
	if err != nil {
		log.Fatal().Err(err).Msg("di.DailyPlayersDataS3Repository failed")
	}

	return dr
}

func NewPriceReportFsRepository() *fs.PriceReportRepository {
	rr, err := fs.NewPriceReportRepository(filepath.Join(fsDir(), "reports"))
	if err != nil {
		log.Fatal().Err(err).Msg("di.NewPriceReportFsRepository failed")
	}

	return rr
}

func NewPriceReportS3Repository() *s3.PriceReportRepository {
	cfg := Config()

	rr, err := s3.NewPriceReportRepository(cfg.AWS, "reports")
	if err != nil {
		log.Fatal().Err(err).Msg("di.NewPriceReportS3Repository failed")
	}

	return rr
}

func NewTeamInMemoryRepository() *memory.TeamRepository {
	tr := memory.NewTeamRepository()

	// NOTE: teams data can be fetched here, since it is quite static.
	wr := Wrapper()

	teams, err := wr.GetTeams()
	if err != nil {
		log.Fatal().Err(err).Msg("di.NewTeamInMemoryRepository - failed to get teams data")
	}

	for _, team := range teams {
		tr.Add(domain.Team(team))
	}

	return &tr
}

// NewFetchServiceCLI builds 'fetch' service for use in CLI build
func NewFetchServiceCLI() *fetch.FetchService {
	wr := Wrapper()
	dr := DailyPlayersDataFsRepository()

	return fetch.NewFetchService(wr, dr)
}

// NewFetchServiceLambdas builds 'fetch' service for use in Lambdas build
func NewFetchServiceLambdas() *fetch.FetchService {
	wr := Wrapper()
	dr := DailyPlayersDataS3Repository()

	return fetch.NewFetchService(wr, dr)

}

// NewGenerateServiceCLI builds 'generate' service for use in CLI build
func NewGenerateServiceCLI() *generate.GenerateService {
	pr := DailyPlayersDataFsRepository()
	rr := NewPriceReportFsRepository()
	tr := NewTeamInMemoryRepository()

	return generate.NewGenerateService(pr, rr, tr)
}

// NewGenerateServiceLambdas builds 'generate' service for use in Lambdas build
func NewGenerateServiceLambdas() *generate.GenerateService {
	pr := DailyPlayersDataS3Repository()
	rr := NewPriceReportS3Repository()
	tr := NewTeamInMemoryRepository()

	return generate.NewGenerateService(pr, rr, tr)
}

// NewGetServiceCLI builds 'get' service for use in CLI build
func NewGetServiceCLI() *get.GetService {
	rr := NewPriceReportFsRepository()

	return get.NewGetService(rr)
}

// NewGetServiceLambdas builds 'get' service for use in Lambdas build
func NewGetServiceLambdas() *get.GetService {
	rr := NewPriceReportS3Repository()

	return get.NewGetService(rr)
}

func NewServer() *echo.Echo {
	gs := NewGetServiceCLI()

	s := rest.NewServer(gs)

	return s
}
