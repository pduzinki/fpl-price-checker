package di

import (
	"os/user"
	"path/filepath"

	"github.com/pduzinki/fpl-price-checker/pkg/config"
	"github.com/pduzinki/fpl-price-checker/pkg/rest"
	"github.com/pduzinki/fpl-price-checker/pkg/services/fetch"
	"github.com/pduzinki/fpl-price-checker/pkg/services/generate"
	"github.com/pduzinki/fpl-price-checker/pkg/services/get"
	"github.com/pduzinki/fpl-price-checker/pkg/storage/fs"
	"github.com/pduzinki/fpl-price-checker/pkg/storage/s3"
	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"

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

// TODO DRY: merge fs and s3 service constructors

func NewPriceReportS3Repository() *fs.PriceReportRepository {
	rr, err := fs.NewPriceReportRepository("reports")
	if err != nil {
		log.Fatal().Err(err).Msg("di.NewPriceReportS3Repository failed")
	}

	return rr
}

func NewFetchService() *fetch.FetchService {
	wr := Wrapper()
	dr := DailyPlayersDataFsRepository()

	return fetch.NewFetchService(wr, dr)
}

func NewFetchServiceS3() *fetch.FetchService {
	wr := Wrapper()
	dr := DailyPlayersDataS3Repository()

	return fetch.NewFetchService(wr, dr)

}

func NewGenerateService() *generate.GenerateService {
	pr := DailyPlayersDataFsRepository()
	rr := NewPriceReportFsRepository()

	return generate.NewGenerateService(pr, rr)
}

func NewGenerateServiceS3() *generate.GenerateService {
	pr := DailyPlayersDataFsRepository()
	rr := NewPriceReportS3Repository()

	return generate.NewGenerateService(pr, rr)
}

func NewGetService() *get.GetService {
	rr := NewPriceReportFsRepository()

	return get.NewGetService(rr)
}

func NewGetServiceS3() *get.GetService {
	rr := NewPriceReportS3Repository()

	return get.NewGetService(rr)
}

func NewServer() *echo.Echo {
	gs := NewGetService()

	s := rest.NewServer(gs)

	return s
}
