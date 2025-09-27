package utils

import (
	"time"
	_ "time/tzdata"
)

var TimeUtil = newTimeUtil("Asia/Shanghai")
var TimeUtilUTC = newTimeUtil(time.UTC.String())

type timeUtil struct {
	Zone     string
	Location *time.Location
}

func newTimeUtil(zone string) *timeUtil {
	location, _ := time.LoadLocation(zone)
	return &timeUtil{
		Zone:     zone,
		Location: location,
	}
}

func (*timeUtil) Now() int64 {
	return time.Now().UnixMilli()
}

func (*timeUtil) Microtime() int64 {
	return time.Now().UnixMicro()
}

func (util *timeUtil) ToString(t time.Time) string {
	return t.In(util.Location).Format("2006-01-02 15:04:05")
}

func (util *timeUtil) TodayStartTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, util.Location)
}

func (util *timeUtil) TodayStartTimeByTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, util.Location)
}

func (util *timeUtil) WeekStartTime() time.Time {
	return util.WeekStartTimeByTime(time.Now())
}

func (util *timeUtil) WeekStartTimeByTime(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := t.AddDate(0, 0, -weekday+1)
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, util.Location)
}

func (util *timeUtil) MonthStartTime() time.Time {
	return util.MonthStartTimeByTime(time.Now())
}

func (util *timeUtil) MonthStartTimeByTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, util.Location)
}
