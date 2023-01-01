package di

import (
	"fmt"

	"github.com/pduzinki/fpl-price-checker/pkg/services/fetch"
	"github.com/pduzinki/fpl-price-checker/pkg/services/generate"
	"github.com/pduzinki/fpl-price-checker/pkg/storage/fs"
	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"
)

func NewFetchService() (*fetch.FetchService, error) {
	wr := wrapper.NewWrapper()

	st, err := fs.NewDailyPlayersDataRepository("./data")
	if err != nil {
		return nil, fmt.Errorf("di NewFetchService, failed to create players data repository: %w", err)
	}

	return fetch.NewFetchService(&wr, &st), nil

}

func NewGenerateService() (*generate.GenerateService, error) {
	st, err := fs.NewDailyPlayersDataRepository("./data")
	if err != nil {
		return nil, fmt.Errorf("di NewGenerateService, failed to create players data repository: %w", err)
	}

	return generate.NewGenerateService(&st), nil
}
