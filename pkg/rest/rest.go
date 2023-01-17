package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/rs/zerolog/log"
)

type ReportGetter interface {
	GetLatestReport(ctx context.Context) (domain.PriceChangeReport, error)
	GetReportByDate(ctx context.Context, date string) (domain.PriceChangeReport, error)
}

func NewServer(rg ReportGetter) *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/latest")
	})
	e.GET("/latest", GetLatest(rg))
	e.GET("/:date", GetByDate(rg))

	return e
}

func GetLatest(rs ReportGetter) func(c echo.Context) error {
	return func(c echo.Context) error {

		report, err := rs.GetLatestReport(c.Request().Context())
		if err != nil {
			log.Error().Err(err).Msg("rest.GetLatest failed to get report")

			return err
		}

		c.JSONPretty(http.StatusOK, report, "  ")

		return nil
	}
}

func GetByDate(rs ReportGetter) func(c echo.Context) error {
	return func(c echo.Context) error {
		date := c.Param("date")

		report, err := rs.GetReportByDate(c.Request().Context(), date)
		if err != nil {
			log.Error().Err(err).Msg("rest.GetByDate failed to get report")

			return err
		}

		c.JSONPretty(http.StatusOK, report, "  ")

		return nil
	}
}
