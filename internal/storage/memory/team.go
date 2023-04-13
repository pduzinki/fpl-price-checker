package memory

import (
	"sync"

	"github.com/pduzinki/fpl-price-checker/internal/domain"
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

func (tr *TeamRepository) GetAll() (map[int]domain.Team, error) {
	tr.RLock()
	defer tr.RUnlock()

	teams := make(map[int]domain.Team)

	for _, team := range tr.teams {
		teams[team.ID] = domain.Team(team)
	}

	return teams, nil
}
