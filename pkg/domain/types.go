package domain

import "time"

const DateFormat = "2006-01-02"

type Player struct {
	ID         int
	Name       string
	Price      int
	SelectedBy string
}

type Record struct {
	Name        string
	OldPrice    string
	NewPrice    string
	Description string
}

type PriceChangeReport struct {
	Date    string
	Records []Record
}

func ParseDate(date string) error {
	_, err := time.Parse(DateFormat, date)

	return err
}
