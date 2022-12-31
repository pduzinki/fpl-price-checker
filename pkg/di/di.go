package di

import (
	"github.com/pduzinki/fpl-price-checker/pkg/services/fetch"
	"github.com/pduzinki/fpl-price-checker/pkg/storage/fs"
	"github.com/pduzinki/fpl-price-checker/pkg/wrapper"
)

func NewFetchService() fetch.FetchService {
	wr := wrapper.NewWrapper()

	st, err := fs.NewDailyPlayersDataRepository("./data")
	if err != nil {
		panic(err)
	}

	return fetch.NewFetchService(&wr, &st)

}
