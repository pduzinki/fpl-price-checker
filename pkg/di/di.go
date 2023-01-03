package di

import (
	"github.com/pduzinki/fpl-price-checker/pkg/rest"
	"github.com/pduzinki/fpl-price-checker/pkg/services/fetch"
	"github.com/pduzinki/fpl-price-checker/pkg/services/generate"
	"github.com/pduzinki/fpl-price-checker/pkg/storage/fs"
	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func Wrapper() *wrapper.Wrapper {
	return wrapper.NewWrapper()
}

func DailyPlayersDataRepository() *fs.DailyPlayersDataRepository {
	dr, err := fs.NewDailyPlayersDataRepository("./data/players/")
	if err != nil {
		log.Fatal().Err(err).Msg("di.DailyPlayersDataRepository failed")
	}

	return dr
}

func NewPriceReportRepository() *fs.PriceReportRepository {
	rr, err := fs.NewPriceReportRepository("./data/reports/")
	if err != nil {
		log.Fatal().Err(err).Msg("di.NewPriceReportRepository failed")
	}

	return rr
}

func NewFetchService() (*fetch.FetchService, error) {
	wr := Wrapper()
	dr := DailyPlayersDataRepository()

	return fetch.NewFetchService(wr, dr), nil

}

func NewGenerateService() (*generate.GenerateService, error) {
	pr := DailyPlayersDataRepository()
	rr := NewPriceReportRepository()

	return generate.NewGenerateService(pr, rr), nil
}

func NewServer() *echo.Echo {
	return rest.NewServer()
}
