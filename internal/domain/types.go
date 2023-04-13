package domain

import "time"

const DateFormat = "2006-01-02"

type Player struct {
	ID         int
	TeamID     int
	Name       string
	Price      int
	SelectedBy string
}

type DailyPlayersData map[int]Player

type Record struct {
	Name        string
	Team        string
	OldPrice    string
	NewPrice    string
	Description string
}

type PriceChangeReport struct {
	Date    string
	Records []Record
}

type Team struct {
	ID        int
	Name      string
	Shortname string
}

func ParseDate(date string) error {
	_, err := time.Parse(DateFormat, date)

	return err
}
