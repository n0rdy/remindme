package utils

import (
	"n0rdy.me/remindme/common"
	"time"
)

const (
	AM = "AM"
	PM = "PM"
)

func TimeFrom24HoursString(timeAs24HoursString string) (time.Time, error) {
	now := time.Now()

	parsedTime, err := time.Parse(common.TimeFormat24Hours, timeAs24HoursString+":00")
	if err != nil {
		return now, common.ErrCmdWrongFormatted24HoursTime
	}

	parsedTime = time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, time.Local)
	if parsedTime.Before(now) || parsedTime.Equal(now) {
		return now, common.ErrCmdTimeShouldBeInFuture
	}

	return parsedTime, nil
}

func TimeFrom12HoursAmPmString(timeAs12HoursAmPmString string, amOrPm string) (time.Time, error) {
	now := time.Now()

	parsedTime, err := time.Parse(common.TimeFormat12AmPmHours, timeAs12HoursAmPmString+" "+amOrPm)
	if err != nil {
		return now, common.ErrCmdWrongFormatted12HoursAmPmTime
	}

	parsedTime = time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
	if parsedTime.Before(now) || parsedTime.Equal(now) {
		return now, common.ErrCmdTimeShouldBeInFuture
	}

	return parsedTime, nil
}

func AddDuration(t time.Time, seconds int, minutes int, hours int) time.Time {
	return t.Local().Add(
		time.Second*time.Duration(seconds) +
			time.Minute*time.Duration(minutes) +
			time.Hour*time.Duration(hours),
	)
}
