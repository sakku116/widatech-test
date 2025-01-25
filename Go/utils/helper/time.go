package helper

import (
	"time"
)

func TimeNowUTC() time.Time {
	return time.Now().UTC()
}

func TimeNowEpochUtc() int64 {
	return TimeNowUTC().Unix()
}

// layout: DD-MM-YYYY
func ParseDateString(date string) (*time.Time, error) {
	parsedDate, err := time.Parse("02-01-2006", date)
	if err != nil {
		return nil, err
	}

	return &parsedDate, nil
}
