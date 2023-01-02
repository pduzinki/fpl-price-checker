package domain

import "time"

const DateFormat = "2006-01-02"

type Record struct {
}

type PriceChangeReport struct {
	Records []Record
}

func ParseDate(date string) error {
	_, err := time.Parse(DateFormat, date)

	return err
}
