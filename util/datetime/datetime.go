package datetime

import (
	"fmt"
	"log"
	"time"
)

const (
	DAY   = 24 * time.Hour
	WEEK  = 7 * DAY
	MONTH = 31 * DAY
	YEAR  = 365 * DAY
)

func NowInEasternCanada() time.Time {
	now := time.Now()
	loc, err := time.LoadLocation("Canada/Eastern")
	if err != nil {
		log.Panicf("unable to load location: %v", err)
	}
	now = now.In(loc)

	// return time without time zone
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
}

func DurationString(date time.Time) string {
	now := NowInEasternCanada()

	if date.After(now) {
		return "INVALID"
	}

	duration := now.Sub(date)

	if duration < time.Minute {
		return DurationMessage(int(duration.Seconds()), "second")
	} else if duration < time.Hour {
		return DurationMessage(int(duration.Minutes()), "minute")
	} else if duration < DAY {
		return DurationMessage(int(duration.Hours()), "hour")
	} else if duration < WEEK {
		return DurationMessage(int(duration.Hours()/24), "day")
	} else if duration < MONTH {
		return DurationMessage(int(duration.Hours()/168), "week")
	} else if duration < YEAR {
		return DurationMessage(int(duration.Hours()/730), "month")
	}

	return "more than a year ago"
}

func DurationMessage(duration int, message string) string {
	if duration == 1 {
		return fmt.Sprintf("1 %s ago", message)
	}

	return fmt.Sprintf("%d %ss ago", duration, message)
}
