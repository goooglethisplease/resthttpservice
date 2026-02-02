package entity

import (
	"fmt"
	"time"
)

const monthYear = "01-2006"

func ParseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse(monthYear, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q, expected MM-YYYY: %w", s, err)
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func FormatMonthYear(t time.Time) string {
	return t.UTC().Format(monthYear)
}
