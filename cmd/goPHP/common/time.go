package common

import "time"

func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func IsLeapYear(year int) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}

func Iso8601Weekday(weekday time.Weekday) int {
	switch weekday {
	case time.Sunday:
		return 7
	default:
		return int(weekday)
	}
}
