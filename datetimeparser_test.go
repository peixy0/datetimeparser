package datetimeparser

import (
	"testing"
	"time"
)

func assert(t *testing.T, value any, expected any, msg string) {
	if value != expected {
		t.Errorf("%v %v != %v (expected)", msg, value, expected)
	}
}

func TestParseStandardFullFormat(t *testing.T) {
	dateParser := NewDateTimeParser(time.Now())
	r, err := dateParser.ParseDateTime("2016年8月12日下午3点14分")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2016, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 12, "day mismatch")
	assert(t, r.Hour(), 15, "hour mismatch")
	assert(t, r.Minute(), 14, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseStandardFullFormatOmitMinuteUnit(t *testing.T) {
	dateParser := NewDateTimeParser(time.Now())
	r, err := dateParser.ParseDateTime("2016年11月12日下午3点14")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2016, "year mismatch")
	assert(t, r.Month(), time.November, "month mismatch")
	assert(t, r.Day(), 12, "day mismatch")
	assert(t, r.Hour(), 15, "hour mismatch")
	assert(t, r.Minute(), 14, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseTomorrowMorning(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2006, time.January, 13, 13, 45, 55, 12, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("明天上午8点")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2006, "year mismatch")
	assert(t, r.Month(), time.January, "month mismatch")
	assert(t, r.Day(), 14, "day mismatch")
	assert(t, r.Hour(), 8, "hour mismatch")
	assert(t, r.Minute(), 0, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseDayAfterTomorrowEvening(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2006, time.December, 31, 13, 45, 55, 12, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("后天晚上11点半")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2007, "year mismatch")
	assert(t, r.Month(), time.January, "month mismatch")
	assert(t, r.Day(), 2, "day mismatch")
	assert(t, r.Hour(), 23, "hour mismatch")
	assert(t, r.Minute(), 30, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseTimeWithQuarter(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2006, time.December, 31, 13, 45, 55, 12, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("晚上10点一刻")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2006, "year mismatch")
	assert(t, r.Month(), time.December, "month mismatch")
	assert(t, r.Day(), 31, "day mismatch")
	assert(t, r.Hour(), 22, "hour mismatch")
	assert(t, r.Minute(), 15, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseThisMonth(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 20, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("这个月31号早上8点15")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 31, "day mismatch")
	assert(t, r.Hour(), 8, "hour mismatch")
	assert(t, r.Minute(), 15, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseThisWeek(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 20, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("周一早上三点三刻")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 15, "day mismatch")
	assert(t, r.Hour(), 3, "hour mismatch")
	assert(t, r.Minute(), 45, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseNextWeek(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 20, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("下周一下午三点半")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 22, "day mismatch")
	assert(t, r.Hour(), 15, "hour mismatch")
	assert(t, r.Minute(), 30, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseThisSunday(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 20, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("周日下午三点半")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 21, "day mismatch")
	assert(t, r.Hour(), 15, "hour mismatch")
	assert(t, r.Minute(), 30, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseThisSunday2(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 21, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("周日下午三点半")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 21, "day mismatch")
	assert(t, r.Hour(), 15, "hour mismatch")
	assert(t, r.Minute(), 30, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseSkipSunday(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 12, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("下周日早上三点三刻")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 21, "day mismatch")
	assert(t, r.Hour(), 3, "hour mismatch")
	assert(t, r.Minute(), 45, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseMinutePeriod(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 12, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("一分钟后")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 12, "day mismatch")
	assert(t, r.Hour(), 12, "hour mismatch")
	assert(t, r.Minute(), 35, "minute mismatch")
	assert(t, r.Second(), 56, "second mismatch")
}

func TestParseHourMinutePeriod(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 12, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("两小时三分钟后")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 12, "day mismatch")
	assert(t, r.Hour(), 14, "hour mismatch")
	assert(t, r.Minute(), 37, "minute mismatch")
	assert(t, r.Second(), 56, "second mismatch")
}

func TestParseOmitDate(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 12, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("15点14分")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 12, "day mismatch")
	assert(t, r.Hour(), 15, "hour mismatch")
	assert(t, r.Minute(), 14, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseYesterday(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 1, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDate("昨天")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.July, "month mismatch")
	assert(t, r.Day(), 31, "day mismatch")
	assert(t, r.Hour(), 0, "hour mismatch")
	assert(t, r.Minute(), 0, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseDayBeforeYesterday(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.January, 1, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("前天早上8点")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2021, "year mismatch")
	assert(t, r.Month(), time.December, "month mismatch")
	assert(t, r.Day(), 30, "day mismatch")
	assert(t, r.Hour(), 8, "hour mismatch")
	assert(t, r.Minute(), 0, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseLastWeekday(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 20, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDate("上周日")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 14, "day mismatch")
	assert(t, r.Hour(), 0, "hour mismatch")
	assert(t, r.Minute(), 0, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseLastWeekday2(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 21, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDate("上周日")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2022, "year mismatch")
	assert(t, r.Month(), time.August, "month mismatch")
	assert(t, r.Day(), 14, "day mismatch")
	assert(t, r.Hour(), 0, "hour mismatch")
	assert(t, r.Minute(), 0, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}

func TestParseLastYear(t *testing.T) {
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	base := time.Date(2022, time.August, 21, 12, 34, 56, 32, shanghai)
	dateParser := NewDateTimeParser(base)
	r, err := dateParser.ParseDateTime("去年1月13日上午8点24分")
	assert(t, err, nil, "error")
	assert(t, r.Year(), 2021, "year mismatch")
	assert(t, r.Month(), time.January, "month mismatch")
	assert(t, r.Day(), 13, "day mismatch")
	assert(t, r.Hour(), 8, "hour mismatch")
	assert(t, r.Minute(), 24, "minute mismatch")
	assert(t, r.Second(), 0, "second mismatch")
}
