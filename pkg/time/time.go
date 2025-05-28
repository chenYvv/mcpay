package time

import (
	"time"
)

const (
	TimeFormatDate     = "2006-01-02"
	TimeFormatDatetime = "2006-01-02 15:04:05"
)

func TowDateDays(startTime time.Time, endTime time.Time) int {
	// 计算天数差异
	duration := endTime.Sub(startTime)
	days := int(duration.Hours() / 24)

	// 如果差异为负数，则说明结束时间在开始时间之前，将结果设为 0
	if days < 0 {
		days = 0
	}
	return days
}

func AddDays(start time.Time, days int) time.Time {
	return start.Add(time.Duration(days) * 24 * time.Hour)
}

func GetDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func StrToTime(timeStr string) time.Time {
	//t, _ := time.Parse(time.DateTime, timeStr)
	t, _ := time.ParseInLocation(time.DateTime, timeStr, time.Local)
	return t
}

func CurToday() (end, sta string) {
	sta = time.Now().Format(TimeFormatDate)
	end = time.Now().Format(TimeFormatDate)
	return
}

func CurThreeDay() (end, sta string) {
	end = time.Now().Format(TimeFormatDate)
	sta = time.Now().AddDate(0, 0, -2).Format(TimeFormatDate)
	return
}

func CurHalfMoth() (end, sta string) {
	end = time.Now().Format(TimeFormatDate)
	sta = time.Now().AddDate(0, 0, -14).Format(TimeFormatDate)
	return
}

func CurWeekDay() (end, sta string) {
	end = time.Now().Format(TimeFormatDate)
	sta = time.Now().AddDate(0, 0, -6).Format(TimeFormatDate)
	return
}

func CurMoth() (end, sta string) {
	end = time.Now().Format(TimeFormatDate)
	sta = time.Now().AddDate(0, -1, 0).Format(TimeFormatDate)
	return
}

func GetMilliseconds(date string) int64 {
	loc, _ := time.LoadLocation("Local") // 或者使用 time.Local
	t, _ := time.ParseInLocation(time.DateTime, date, loc)
	milliseconds := t.UnixNano() / 1e6
	return milliseconds
}
