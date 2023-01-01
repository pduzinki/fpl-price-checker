package generate

import (
	"context"
	"fmt"
	"time"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

const dateFormat = "2006-01-02" // TODO perhaps move to domain?

type StorageGetter interface {
	GetByDate(_ context.Context, date string) (map[int]domain.Player, error)
}

type GenerateService struct {
	sg StorageGetter // TODO find a better name
}

func NewGenerateService(sg StorageGetter) *GenerateService {
	return &GenerateService{
		sg: sg,
	}
}

func (gs *GenerateService) GeneratePriceReport() error {
	todaysDate := time.Now().Format(dateFormat)
	yesterdaysDate := time.Now().Add(-24 * time.Hour).Format(dateFormat)

	yesterdayPlayers, err := gs.sg.GetByDate(context.TODO(), yesterdaysDate)
	if err != nil {
		return fmt.Errorf("err 1: %w", err)
	}

	todayPlayers, err := gs.sg.GetByDate(context.TODO(), todaysDate)
	if err != nil {
		return fmt.Errorf("err 2: %w", err)
	}

	for tk, tv := range todayPlayers {
		yv, prs := yesterdayPlayers[tk]
		if !prs {
			// fmt.Println("new:", tv.Name, tv.Price) // TODO no sure if newly added players should be reported
			continue
		}

		if yv.Price != tv.Price {
			fmt.Println(tv.Name, yv.Price, tv.Price)
		}

		// TODO save price changes in storage, add new repo type
	}

	return nil
}
