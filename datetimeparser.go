package datetimeparser

import (
	"errors"
	"regexp"
	"time"
)

type DateTimeParser struct {
	Base time.Time
}

type DateTimeParseResult struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func NewDateTimeParser(base time.Time) *DateTimeParser {
	return &DateTimeParser{
		Base: base,
	}
}

type ParseFunc[T any] func(string, *T) (string, error)

type ParseFuncList[T any] []ParseFunc[T]

func parseAnyOf[T any](fs []ParseFunc[T]) ParseFunc[T] {
	return func(input string, result *T) (string, error) {
		rest := input
		var err error = nil
		for _, f := range fs {
			r := *result
			rest, err = f(input, &r)
			if err == nil {
				*result = r
				return rest, nil
			}
		}
		return rest, errors.New("not parsed any of")
	}
}

func parseAllOf[T any](fs []ParseFunc[T]) ParseFunc[T] {
	return func(input string, result *T) (string, error) {
		rest := input
		r := *result
		var err error = nil
		for _, f := range fs {
			rest, err = f(rest, &r)
			if err != nil {
				return input, err
			}
		}
		*result = r
		return rest, nil
	}
}

func parseRegex(input string, ex string) (string, error) {
	r, _ := regexp.Compile("^" + ex)
	result := r.Find([]byte(input))
	if result == nil {
		return input, errors.New("regexp not parsed")
	}
	return input[len(result):], nil
}

func parseNumericNumber(input string, r *int) (string, error) {
	n := 0
	parsed := 0
	for parsed < len(input) {
		if input[parsed] >= '0' && input[parsed] <= '9' {
			n *= 10
			n += int(input[parsed] - '0')
			parsed++
			continue
		}
		break
	}
	if parsed == 0 {
		return input, errors.New("number not parsed")
	}
	*r = n
	return input[parsed:], nil
}

func parseNumberWithUnit(input string, unit string, r *int) (string, error) {
	rest, err := parseAnyNumber(input, r)
	if err != nil {
		return rest, err
	}
	rest, err = parseRegex(rest, unit)
	if err != nil {
		return rest, errors.New("expecting unit " + unit)
	}
	return rest, nil
}

var chineseDigits = []string{"(〇|零)", "一", "(二|两)", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}

func parseChineseNumber(input string, r *int) (string, error) {
	for i, d := range chineseDigits {
		rest, err := parseRegex(input, d)
		if err == nil {
			*r = i
			return rest, nil
		}
	}
	return input, errors.New("chinese number not parsed")
}

func parseAnyNumber(input string, r *int) (string, error) {
	return parseAnyOf(ParseFuncList[int]{
		parseNumericNumber,
		parseChineseNumber,
	})(input, r)
}

func parseAnyMinute(input string, r *int) (string, error) {
	rest, err := parseRegex(input, "半")
	if err == nil {
		*r = 30
		return rest, nil
	}
	var k int
	rest, err = parseNumberWithUnit(input, "刻", &k)
	if err == nil {
		*r = k * 15
		return rest, nil
	}
	rest, err = parseNumberWithUnit(input, "(分)?", &k)
	if err == nil {
		*r = k
		return rest, nil
	}
	return input, errors.New("minute not parsed")
}

func parseWeekday(input string, r *int) (string, error) {
	rest, err := parseRegex(input, "(周|星期|礼拜)")
	if err != nil {
		return rest, err
	}
	rest, err = parseRegex(rest, "(天|日)")
	if err == nil {
		*r = 0
		return rest, nil
	}
	var w int
	rest, err = parseChineseNumber(rest, &w)
	if err != nil {
		return rest, errors.New("weekday not parsed")
	}
	if w < 1 || w > 6 {
		return input, errors.New("weekday not parsed")
	}
	*r = w
	return rest, nil
}

func (dp *DateTimeParser) ignore(input string, _ *DateTimeParseResult) (string, error) {
	return input, nil
}

func (dp *DateTimeParser) parseWithHalfHourPeriod(input string, result *DateTimeParseResult) (string, error) {
	var h int = 0
	rest, err := parseAnyOf(ParseFuncList[DateTimeParseResult]{
		func(input string, _ *DateTimeParseResult) (string, error) {
			return parseNumberWithUnit(input, "个半(小时|钟头)(以)?后", &h)
		},
		func(input string, _ *DateTimeParseResult) (string, error) {
			return parseRegex(input, "半(个)?(小时|钟头)(以)?后")
		},
	})(input, result)
	if err != nil {
		return input, err
	}
	t := dp.Base.Add(time.Duration(h)*time.Hour + 30*time.Minute)
	result.Hour = t.Hour()
	result.Minute = t.Minute()
	result.Second = dp.Base.Second()
	return rest, nil
}

func (dp *DateTimeParser) parseHourPeriod(input string, result *DateTimeParseResult) (string, error) {
	var h int
	rest, err := parseNumberWithUnit(input, "(个)?(小时|钟头)(以)?后", &h)
	if err != nil {
		return input, err
	}
	t := dp.Base.Add(time.Duration(h) * time.Hour)
	result.Hour = t.Hour()
	result.Minute = t.Minute()
	result.Second = dp.Base.Second()
	return rest, nil
}

func (dp *DateTimeParser) parseMinutePeriod(input string, result *DateTimeParseResult) (string, error) {
	var m int
	rest, err := parseNumberWithUnit(input, "(分钟|分)(以)?后", &m)
	if err != nil {
		return input, err
	}
	t := dp.Base.Add(time.Duration(m) * time.Minute)
	result.Hour = t.Hour()
	result.Minute = t.Minute()
	result.Second = dp.Base.Second()
	return rest, nil
}

func (dp *DateTimeParser) parseHourMinutePeriod(input string, result *DateTimeParseResult) (string, error) {
	var h, m int
	rest, err := parseNumberWithUnit(input, "(个)?(小时|时|钟头)", &h)
	if err != nil {
		return input, err
	}
	rest, err = parseNumberWithUnit(rest, "(分钟|分)(以)?后", &m)
	if err != nil {
		return input, err
	}
	t := dp.Base.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
	result.Hour = t.Hour()
	result.Minute = t.Minute()
	result.Second = dp.Base.Second()
	return rest, nil
}

func (dp *DateTimeParser) parseTimePeriod(input string, result *DateTimeParseResult) (string, error) {
	return parseAnyOf(ParseFuncList[DateTimeParseResult]{
		dp.parseWithHalfHourPeriod,
		dp.parseHourMinutePeriod,
		dp.parseHourPeriod,
		dp.parseMinutePeriod,
	})(input, result)
}

func (dp *DateTimeParser) parseYear(input string, result *DateTimeParseResult) (string, error) {
	var y int
	rest, err := parseNumberWithUnit(input, "年", &y)
	if err != nil {
		return input, err
	}
	result.Year = y
	return rest, nil
}

func (dp *DateTimeParser) parseMonth(input string, result *DateTimeParseResult) (string, error) {
	var m int
	rest, err := parseNumberWithUnit(input, "月", &m)
	if err != nil {
		return input, err
	}
	result.Month = m
	return rest, nil
}

func (dp *DateTimeParser) parseDay(input string, result *DateTimeParseResult) (string, error) {
	var d int
	rest, err := parseNumberWithUnit(input, "(日|号)", &d)
	if err != nil {
		return input, err
	}
	result.Day = d
	return rest, nil
}

func (dp *DateTimeParser) parseYMD(input string, result *DateTimeParseResult) (string, error) {
	return parseAllOf(ParseFuncList[DateTimeParseResult]{dp.parseYear, dp.parseMonth, dp.parseDay})(input, result)
}

func (dp *DateTimeParser) parseMD(input string, result *DateTimeParseResult) (string, error) {
	return parseAllOf(ParseFuncList[DateTimeParseResult]{dp.parseMonth, dp.parseDay})(input, result)
}

func (dp *DateTimeParser) parseLastYear(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "去年")
	if err != nil {
		return rest, err
	}
	n := dp.Base.AddDate(-1, 0, 0)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseNextYear(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "明年")
	if err != nil {
		return rest, err
	}
	n := dp.Base.AddDate(1, 0, 0)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseThisMonth(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "(这(个)?|本)月")
	if err != nil {
		return rest, err
	}
	n := dp.Base
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseLastMonth(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "上个月")
	if err != nil {
		return rest, err
	}
	n := dp.Base.AddDate(0, -1, 0)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseNextMonth(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "下个月")
	if err != nil {
		return rest, err
	}
	n := dp.Base.AddDate(0, 1, 0)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseYesterday(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "昨(天|日)")
	if err != nil {
		return rest, err
	}
	n := dp.Base.AddDate(0, 0, -1)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseDayBeforeYesterday(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "前(天|日)")
	if err != nil {
		return rest, err
	}
	n := dp.Base.AddDate(0, 0, -2)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseToday(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "今(天|日)")
	if err != nil {
		return rest, err
	}
	result.Year = dp.Base.Year()
	result.Month = int(dp.Base.Month())
	result.Day = dp.Base.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseNextDay(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "明(天|日)")
	if err != nil {
		return rest, err
	}
	d := 1
	if dp.Base.Hour() < 5 {
		d = 0
	}
	n := dp.Base.AddDate(0, 0, d)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseDayAfterNextDay(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "后(天|日)")
	if err != nil {
		return rest, err
	}
	d := 2
	if dp.Base.Hour() < 5 {
		d = 1
	}
	n := dp.Base.AddDate(0, 0, d)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseWeekday(input string, result *DateTimeParseResult) (string, error) {
	var w int
	rest, err := parseWeekday(input, &w)
	if err != nil {
		return rest, err
	}
	if w == 0 && int(dp.Base.Weekday()) != 0 {
		w = 7
	}
	d := w - int(dp.Base.Weekday())
	n := dp.Base.AddDate(0, 0, d)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseLastWeekday(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "上")
	if err != nil {
		return rest, err
	}
	var w int
	rest, err = parseWeekday(rest, &w)
	if err != nil {
		return rest, err
	}
	if w == 0 && int(dp.Base.Weekday()) != 0 {
		w = 7
	}
	d := -7 + (w - int(dp.Base.Weekday()))
	n := dp.Base.AddDate(0, 0, d)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseNextWeekday(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "下")
	if err != nil {
		return rest, err
	}
	var w int
	rest, err = parseWeekday(rest, &w)
	if err != nil {
		return rest, err
	}
	w += 7
	d := w - int(dp.Base.Weekday())
	if w == 7 && d < 7 {
		d += 7
	}
	n := dp.Base.AddDate(0, 0, d)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseWeekAfterNextWeekday(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseRegex(input, "下下")
	if err != nil {
		return rest, err
	}
	var w int
	rest, err = parseWeekday(rest, &w)
	if err != nil {
		return rest, err
	}
	w += 7
	d := w - int(dp.Base.Weekday())
	if w == 7 && d < 7 {
		d += 7
	}
	n := dp.Base.AddDate(0, 0, d)
	result.Year = n.Year()
	result.Month = int(n.Month())
	result.Day = n.Day()
	return rest, nil
}

func (dp *DateTimeParser) parseNormHourMinute(input string, result *DateTimeParseResult) (string, error) {
	var h, m int
	rest, err := parseNumericNumber(input, &h)
	if err != nil {
		return rest, err
	}
	rest, err = parseRegex(rest, ":")
	if err != nil {
		return rest, err
	}
	rest, err = parseNumericNumber(rest, &m)
	if err != nil {
		return rest, err
	}
	result.Hour = h
	result.Minute = m
	return rest, nil
}

func (dp *DateTimeParser) parseAmHourMinute(input string, result *DateTimeParseResult) (string, error) {
	return parseAllOf(ParseFuncList[DateTimeParseResult]{
		func(input string, _ *DateTimeParseResult) (string, error) {
			return parseRegex(input, "(上午|凌晨|早上)")
		},
		dp.parseClockTime,
	})(input, result)
}

func (dp *DateTimeParser) parsePmHourMinute(input string, result *DateTimeParseResult) (string, error) {
	rest, err := parseAllOf(ParseFuncList[DateTimeParseResult]{
		func(input string, _ *DateTimeParseResult) (string, error) {
			return parseRegex(input, "(下午|晚上)")
		},
		dp.parseClockTime,
	})(input, result)
	if err == nil && result.Hour < 12 {
		result.Hour += 12
	}
	return rest, err
}

func (dp *DateTimeParser) parseNumberHour(input string, result *DateTimeParseResult) (string, error) {
	var h int
	rest, err := parseNumberWithUnit(input, "(点|时)", &h)
	if err != nil {
		return input, err
	}
	result.Hour = h
	result.Minute = 0
	return rest, nil
}

func (dp *DateTimeParser) parseNumberMinute(input string, result *DateTimeParseResult) (string, error) {
	var m int
	rest, err := parseNumberWithUnit(input, "(分)", &m)
	if err != nil {
		return input, err
	}
	result.Minute = m
	return rest, nil
}

func (dp *DateTimeParser) parseHourMinute(input string, result *DateTimeParseResult) (string, error) {
	var h, m int
	rest, err := parseNumberWithUnit(input, "(点|时)", &h)
	if err != nil {
		return input, err
	}
	rest, err = parseAnyMinute(rest, &m)
	if err != nil {
		return input, err
	}
	result.Hour = h
	result.Minute = m
	return rest, nil
}

func (dp *DateTimeParser) parseAnyDate(input string, result *DateTimeParseResult) (string, error) {
	return parseAnyOf(ParseFuncList[DateTimeParseResult]{
		dp.parseToday,
		dp.parseYesterday,
		dp.parseDayBeforeYesterday,
		dp.parseNextDay,
		dp.parseDayAfterNextDay,
		dp.parseWeekday,
		dp.parseLastWeekday,
		dp.parseNextWeekday,
		dp.parseWeekAfterNextWeekday,
		parseAllOf(ParseFuncList[DateTimeParseResult]{dp.parseThisMonth, dp.parseDay}),
		parseAllOf(ParseFuncList[DateTimeParseResult]{dp.parseLastMonth, dp.parseDay}),
		dp.parseLastMonth,
		parseAllOf(ParseFuncList[DateTimeParseResult]{dp.parseNextMonth, dp.parseDay}),
		dp.parseNextMonth,
		parseAllOf(ParseFuncList[DateTimeParseResult]{dp.parseLastYear, dp.parseMD}),
		dp.parseLastYear,
		parseAllOf(ParseFuncList[DateTimeParseResult]{dp.parseNextYear, dp.parseMD}),
		dp.parseNextYear,
		dp.parseYMD,
		dp.parseMD,
	})(input, result)
}

func (dp *DateTimeParser) parseClockTime(input string, result *DateTimeParseResult) (string, error) {
	return parseAnyOf(ParseFuncList[DateTimeParseResult]{
		dp.parseNormHourMinute,
		dp.parseHourMinute,
		dp.parseNumberHour,
	})(input, result)
}

func (dp *DateTimeParser) parseAnyTime(input string, result *DateTimeParseResult) (string, error) {
	return parseAnyOf(ParseFuncList[DateTimeParseResult]{
		dp.parseAmHourMinute,
		dp.parsePmHourMinute,
		dp.parseClockTime,
	})(input, result)
}

func (dp *DateTimeParser) parseAnyDateTime(input string, result *DateTimeParseResult) (string, error) {
	return parseAnyOf(ParseFuncList[DateTimeParseResult]{
		parseAllOf(ParseFuncList[DateTimeParseResult]{
			dp.parseAnyDate,
			dp.parseAnyTime,
		}),
		dp.parseAnyTime,
	})(input, result)
}

func (dp *DateTimeParser) ParseDateTime(input string) (time.Time, error) {
	result := DateTimeParseResult{
		Year:   dp.Base.Year(),
		Month:  int(dp.Base.Month()),
		Day:    dp.Base.Day(),
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	_, err := parseAnyOf(ParseFuncList[DateTimeParseResult]{
		dp.parseTimePeriod,
		dp.parseAnyDateTime,
	})(input, &result)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(result.Year, time.Month(result.Month), result.Day, result.Hour, result.Minute, result.Second, 0, dp.Base.Location()), nil
}

func (dp *DateTimeParser) ParseDate(input string) (time.Time, error) {
	result := DateTimeParseResult{
		Year:   dp.Base.Year(),
		Month:  int(dp.Base.Month()),
		Day:    dp.Base.Day(),
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	_, err := parseAnyOf(ParseFuncList[DateTimeParseResult]{
		dp.parseAnyDate,
	})(input, &result)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(result.Year, time.Month(result.Month), result.Day, result.Hour, result.Minute, result.Second, 0, dp.Base.Location()), nil
}
