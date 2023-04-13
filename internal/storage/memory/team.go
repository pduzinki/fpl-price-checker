package memory

import (
	"sync"

	"github.com/pduzinki/fpl-price-checker/internal/domain"
	"github.com/pduzinki/fpl-price-checker/internal/storage"
)

type Team struct {
	ID        int
	Name      string
	Shortname string
}

type TeamRepository struct {
	teams map[int]Team
	sync.RWMutex
}

func NewTeamRepository() TeamRepository {
	return TeamRepository{
		teams: make(map[int]Team),
	}
}

func (tr *TeamRepository) Add(teams ...domain.Team) {
	tr.Lock()
	defer tr.Unlock()

	for _, team := range teams {
		tr.teams[team.ID] = Team(team)
	}
}

func (tr *TeamRepository) GetByID(id int) (domain.Team, error) {
	tr.RLock()
	defer tr.RUnlock()

	if team, prs := tr.teams[id]; prs {
		return domain.Team(team), nil
	}

	return domain.Team{}, storage.ErrDataNotFound
}
