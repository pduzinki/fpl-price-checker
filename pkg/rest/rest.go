package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"go.uber.org/multierr"

	"github.com/rs/zerolog/log"
)

type ReportGetter interface {
	GetByDate(ctx context.Context, date string) (domain.PriceChangeReport, error)
}

func NewServer(rg ReportGetter) *echo.Echo {
	e := echo.New()

	e.GET("/latest", GetLatest(rg))
	e.GET("/:date", GetByDate(rg))

	return e
}

func GetLatest(rg ReportGetter) func(c echo.Context) error {
	return func(c echo.Context) error {
		todaysDate := time.Now().Format(domain.DateFormat)
		yesterdaysDate := time.Now().Add(-24 * time.Hour).Format(domain.DateFormat)

		var errors error

		// NOTE: FPL price changes usually occur around 1:30am GMT, so there's a time gap,
		// where there's no "today's" report just yet. in that case, latest will be a report from the day before.

		report, err := rg.GetByDate(c.Request().Context(), todaysDate) // TODO use multi err??
		if err == nil {
			c.JSONPretty(http.StatusOK, report, "  ")

			return nil
		}
		errors = multierr.Append(errors, err)

		report, err = rg.GetByDate(c.Request().Context(), yesterdaysDate)
		if err == nil {
			c.JSONPretty(http.StatusOK, report, "  ")

			return nil
		}
		errors = multierr.Append(errors, err)

		log.Error().Err(err).Msg("rest.GetLatest failed to get report")
		return errors
	}
}

func GetByDate(rg ReportGetter) func(c echo.Context) error {
	return func(c echo.Context) error {
		date := c.Param("date")

		if err := domain.ParseDate(date); err != nil {
			return err
		}

		report, err := rg.GetByDate(c.Request().Context(), date)
		if err != nil {
			log.Error().Err(err).Msg("rest.GetByDate failed to get report")
			return err
		}

		c.JSONPretty(http.StatusOK, report, "  ")

		return nil
	}
}
