package utils

import (
	"github.com/lestrrat-go/strftime"
	"strconv"
	"time"
)

func GetFormatTime(timeStamp int64, timeFormat string) string {
	if timeStamp == 0 {
		timeStamp = 1
	}
	f, _ := strftime.New(timeFormat)
	return f.FormatString(time.Unix(timeStamp, 0))
}

func GetTimeObj(timeIn time.Time) map[string]string {
	var timeObj = map[string]string{
		"Year":      strconv.Itoa(timeIn.Year()),
		"Month":     timeIn.Month().String(),
		"Day":       strconv.Itoa(timeIn.Day()),
		"Hour":      strconv.Itoa(timeIn.Day()),
		"Minute":    strconv.Itoa(timeIn.Minute()),
		"Second":    strconv.Itoa(timeIn.Second()),
		"TimeStamp": strconv.FormatInt(timeIn.Unix(), 10),
	}
	return timeObj

}
