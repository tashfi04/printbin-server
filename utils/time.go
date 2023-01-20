package utils

import (
	"time"
)

func DefaultUpdateTime() time.Time {

	t, _ := time.Parse(time.RFC822Z, time.RFC822Z)
	return t
}

func ParseDate(stringDate string) (*time.Time, error) {

	//date, err := time.Parse(time.RFC3339, stringDate)
	date, err := time.Parse(time.RFC3339, stringDate)
	if err != nil {
		return nil, err
	}

	date = date.UTC()

	return &date, nil
}
