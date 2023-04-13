package wrapper

type Bootstrap struct {
	// we're only interested in players data
	Players []Player `json:"elements"`
	// and teams data
	Teams []Team `json:"teams"`
}

type Player struct {
	ID                  int    `json:"id"`
	Team                int    `json:"team"`
	Position            int    `json:"element_type"`
	WebName             string `json:"web_name"`
	Price               int    `json:"now_cost"`
	SelectedBy          string `json:"selected_by_percent"`
	CostChangeEvent     int    `json:"cost_change_event"`
	CostChangeEventFall int    `json:"cost_change_event_fall"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"short_name"`
}
