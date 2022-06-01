package utils

import "time"

func StartOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

func EndOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 24, 0, 0, -1, date.Location())
}
