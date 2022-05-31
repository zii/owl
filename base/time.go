package base

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// 生日计算年龄
// Date format: 2006-01-02
func GetAge(fromDate string) (int, error) {
	st, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return 0, err
	}
	sy, sm, sd := st.Date()
	now := time.Now()
	ny, nm, nd := now.Date()
	if nm < sm || nm == sm && nd < sd {
		return ny - sy - 1, nil
	} else {
		return ny - sy, nil
	}
}

// 时间戳转日期时间
func FromUnixTime(timestamp int) string {
	return time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
}

// 时间戳转日期 2019-07-17
func ToDate(timestamp int) string {
	return time.Unix(int64(timestamp), 0).Format("2006-01-02")
}

// 日期字符串转成时间戳
func FromDate(date string) (int, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}
	t, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return 0, err
	}
	return int(t.Local().Unix()), nil
}

// 时间戳转日期 20190717, 负数表示今天或今天之前的相对秒数
func ToYMD(timestamp int) string {
	if timestamp <= 0 {
		timestamp = int(time.Now().Unix()) + timestamp
	}
	return time.Unix(int64(timestamp), 0).Format("20060102")
}

func ToIntYMD(time time.Time) int {
	y, m, d := time.Date()
	return y*10000 + int(m)*100 + d
}

// 时间戳转日期时间 20190717150405, 负数表示今天或今天之前的相对秒数
func ToYMDHIS(timestamp int) string {
	if timestamp <= 0 {
		timestamp = int(time.Now().Unix()) + timestamp
	}
	return time.Unix(int64(timestamp), 0).Format("20060102150405")
}

func FormatDate(time time.Time, style string) string {
	layout := string(style)
	layout = strings.Replace(layout, "yyyy", "2006", 1)
	layout = strings.Replace(layout, "yy", "06", 1)
	layout = strings.Replace(layout, "MM", "01", 1)
	layout = strings.Replace(layout, "dd", "02", 1)
	layout = strings.Replace(layout, "HH", "15", 1)
	layout = strings.Replace(layout, "mm", "04", 1)
	layout = strings.Replace(layout, "ss", "05", 1)
	layout = strings.Replace(layout, "SSS", "000", -1)

	return time.Format(layout)
}

// 获取(本月1号, 下个月1号)时间戳
func GetMonthRange(now time.Time) (int, int) {
	y, m, _ := now.Date()
	nextm := (m + 1) % 12
	var nexty int
	if nextm < m {
		nexty = y + 1
	} else {
		nexty = y
	}
	sd := time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
	ed := time.Date(nexty, nextm, 1, 0, 0, 0, 0, now.Location())
	return int(sd.Unix()), int(ed.Unix())
}

func IsInSameHour(t1 time.Time, t2 time.Time) bool {
	return GetHourStart(t1) == GetHourStart(t2)
}

func GetHourStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

func IsInSameWeek(t1 time.Time, t2 time.Time) bool {
	return GetWeekStart(t1) == GetWeekStart(t2)
}

//获取某一天的0点时间
func GetDayStart(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func Today() int {
	s := GetDayStart(time.Now())
	return int(s.Unix())
}

// 周1零点
func GetWeekStart(t time.Time) time.Time {
	weekday := t.Weekday()
	diffday := 0
	if weekday == time.Sunday {
		diffday = -6
	} else {
		diffday = -(int(weekday) - 1)
	}

	t2 := t.Add(time.Duration(diffday) * time.Hour * 24)

	return GetDayStart(t2)
}

func GetWeekEnd(t time.Time) time.Time {
	return GetWeekStart(t).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// 月初0点
func GetMonthStart(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetDayStart(d)
}

// 月初0点
func GetMonthEnd(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetDayStart(d).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func GetConstellationMD(month, day int) int {
	if month <= 0 || month >= 13 {
		return 0
	}
	if day <= 0 || day >= 32 {
		return 0
	}
	if (month == 1 && day >= 20) || (month == 2 && day <= 18) {
		return 11
	}
	if (month == 2 && day >= 19) || (month == 3 && day <= 20) {
		return 12
	}
	if (month == 3 && day >= 21) || (month == 4 && day <= 19) {
		return 1
	}
	if (month == 4 && day >= 20) || (month == 5 && day <= 20) {
		return 2
	}
	if (month == 5 && day >= 21) || (month == 6 && day <= 21) {
		return 3
	}
	if (month == 6 && day >= 22) || (month == 7 && day <= 22) {
		return 4
	}
	if (month == 7 && day >= 23) || (month == 8 && day <= 22) {
		return 5
	}
	if (month == 8 && day >= 23) || (month == 9 && day <= 22) {
		return 6
	}
	if (month == 9 && day >= 23) || (month == 10 && day <= 22) {
		return 7
	}
	if (month == 10 && day >= 23) || (month == 11 && day <= 21) {
		return 8
	}
	if (month == 11 && day >= 22) || (month == 12 && day <= 21) {
		return 9
	}
	if (month == 12 && day >= 22) || (month == 1 && day <= 19) {
		return 10
	}
	return 0
}

// 返回星座(1-12)
func GetConstellation(date string) int {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0
	}
	return GetConstellationMD(int(t.Month()), t.Day())
}

func HourIn(fromHour int, endHour int) bool {
	h := time.Now().Hour()
	return h >= fromHour && h <= endHour
}

// 两个时间戳的日期相减差多少天 a-b
func DayDiff(a, b int) int {
	return (a-16*3600)/86400 - (b-16*3600)/86400
}

// 周一ymd
func Monday(t time.Time) string {
	weekday := t.Weekday()
	diffday := 0
	if weekday == time.Sunday {
		diffday = -6
	} else {
		diffday = -(int(weekday) - 1)
	}

	t2 := t.Add(time.Duration(diffday) * time.Hour * 24)
	return t2.Format("20060102")
}

// 今日ymd
func Day(t time.Time) string {
	return t.Format("20060102")
}

// 本月ym
func Month(t time.Time) string {
	return t.Format("200601")
}
