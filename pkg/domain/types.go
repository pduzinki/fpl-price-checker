package domain

import "time"

const DateFormat = "2006-01-02"

type Record struct {
	Name        string
	OldPrice    float64
	NewPrice    float64
	Description string
}

type PriceChangeReport struct {
	Records []Record
}

func ParseDate(date string) error {
	_, err := time.Parse(DateFormat, date)

	return err
}
