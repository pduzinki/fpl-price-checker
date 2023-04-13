package generate

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/pduzinki/fpl-price-checker/internal/domain"
)

//go:generate moq -out generate_moq_test.go . DailyPlayersDataGetter PriceChangeReportAdder TeamsGetter

type Record struct {
	Name        string
	Team        string
	OldPrice    string
	NewPrice    string
	Description string
	SelectedBy  float64
}

type DailyPlayersDataGetter interface {
	GetByDate(ctx context.Context, date string) (domain.DailyPlayersData, error)
}

type PriceChangeReportAdder interface {
	Add(ctx context.Context, date string, report domain.PriceChangeReport) error
}

type TeamsGetter interface {
	GetAll() (map[int]domain.Team, error)
}

type GenerateService struct {
	dg DailyPlayersDataGetter
	ra PriceChangeReportAdder
	tg TeamsGetter
}

func NewGenerateService(sg DailyPlayersDataGetter, ra PriceChangeReportAdder, tg TeamsGetter) *GenerateService {
	return &GenerateService{
		dg: sg,
		ra: ra,
		tg: tg,
	}
}

func (gs *GenerateService) GeneratePriceReport(ctx context.Context) error {
	todaysDate := time.Now().Format(domain.DateFormat)
	yesterdaysDate := time.Now().Add(-24 * time.Hour).Format(domain.DateFormat)

	yesterdayPlayers, err := gs.dg.GetByDate(ctx, yesterdaysDate)
	if err != nil {
		return fmt.Errorf("generate.GenerateService.GeneratePriceReport failed to get yesterday's players data: %w", err)
	}

	todayPlayers, err := gs.dg.GetByDate(ctx, todaysDate)
	if err != nil {
		return fmt.Errorf("generate.GenerateService.GeneratePriceReport failed to get today's players data: %w", err)
	}

	teams, err := gs.tg.GetAll()
	if err != nil {
		return fmt.Errorf("generate.GenerateService.GeneratePriceReport failed to get teams data: %w", err)
	}

	report := domain.PriceChangeReport{
		Date:    todaysDate,
		Records: generateRecords(yesterdayPlayers, todayPlayers, teams),
	}

	err = gs.ra.Add(ctx, todaysDate, report)
	if err != nil {
		return fmt.Errorf("generate.GenerateService.GeneratePriceReport failed to save report: %w", err)
	}

	return nil
}

func generateRecords(yesterdayPlayers, todayPlayers domain.DailyPlayersData, teams map[int]domain.Team) []domain.Record {
	priceChangedPlayers := make([]Record, 0)
	newPlayers := make([]Record, 0)

	for tk, tv := range todayPlayers {
		yv, prs := yesterdayPlayers[tk]
		if !prs {
			newPlayers = append(newPlayers, Record{
				Name:        tv.Name,
				Team:        addTeam(teams, tv.TeamID),
				OldPrice:    "-",
				NewPrice:    fmt.Sprintf("%.1f", float64(tv.Price)/10.),
				Description: "new",
			})
			continue
		}

		if yv.Price != tv.Price {
			selectedBy, err := strconv.ParseFloat(tv.SelectedBy, 64)
			if err != nil {
				selectedBy = 0.0
			}

			record := Record{
				Name:        tv.Name,
				Team:        addTeam(teams, tv.TeamID),
				OldPrice:    fmt.Sprintf("%.1f", float64(yv.Price)/10.),
				NewPrice:    fmt.Sprintf("%.1f", float64(tv.Price)/10.),
				Description: addDescription(yv.Price, tv.Price),
				SelectedBy:  selectedBy,
			}

			priceChangedPlayers = append(priceChangedPlayers, record)
		}
	}

	sort.Slice(priceChangedPlayers, func(i, j int) bool {
		return priceChangedPlayers[i].SelectedBy > priceChangedPlayers[j].SelectedBy
	})

	priceChangedPlayers = append(priceChangedPlayers, newPlayers...)

	records := make([]domain.Record, 0, len(priceChangedPlayers))
	for _, player := range priceChangedPlayers {
		records = append(records, toDomainRecord(player))
	}

	return records
}

func addDescription(oldPrice, newPrice int) string {
	if oldPrice > newPrice {
		return "drop"
	}
	return "rise"
}

func addTeam(teams map[int]domain.Team, teamID int) string {
	if team, prs := teams[teamID]; prs {
		return team.Shortname
	}

	return "-"
}

func toDomainRecord(record Record) domain.Record {
	return domain.Record{
		Name:        record.Name,
		Team:        record.Team,
		OldPrice:    record.OldPrice,
		NewPrice:    record.NewPrice,
		Description: record.Description,
	}
}
