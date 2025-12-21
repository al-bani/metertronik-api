package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type TimeData struct {
	Time time.Time
}

func ToUTC(t time.Time) time.Time {
	return t.In(time.UTC)
}

func (t TimeData) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.UTC(), nil
}

func (t *TimeData) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	var parsedTime time.Time
	var err error

	switch v := value.(type) {
	case time.Time:
		parsedTime = v
	case []byte:
		parsedTime, err = time.Parse(time.RFC3339, string(v))
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", string(v))
			if err != nil {
				parsedTime, err = time.Parse("2006-01-02", string(v))
				if err != nil {
					return err
				}
			}
		}
	case string:
		parsedTime, err = time.Parse(time.RFC3339, v)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", v)
			if err != nil {
				parsedTime, err = time.Parse("2006-01-02", v)
				if err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("cannot scan TimeData from value")
	}

	t.Time = parsedTime.UTC()
	return nil
}

func (t TimeData) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.UTC().Format(time.RFC3339))
}

func (t *TimeData) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}

	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}

	t.Time = parsedTime.UTC()
	return nil
}

func (t TimeData) Format() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}

func (t TimeData) FormatUTC() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}

func (t TimeData) AddHours(hours int) TimeData {
	return TimeData{Time: t.Time.UTC().Add(time.Duration(hours) * time.Hour)}
}

func (t TimeData) AddMinutes(minutes int) TimeData {
	return TimeData{Time: t.Time.UTC().Add(time.Duration(minutes) * time.Minute)}
}

func (t TimeData) AddDays(days int) TimeData {
	return TimeData{Time: t.Time.UTC().Add(time.Duration(days) * 24 * time.Hour)}
}

func Minutes(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}

func Hours(hours int) time.Duration {
	return time.Duration(hours) * time.Hour
}

func Days(days int) time.Duration {
	return time.Duration(days) * 24 * time.Hour
}

func (t TimeData) StartOfDay() TimeData {
	utcTime := t.Time.UTC()
	year, month, day := utcTime.Date()

	startUTC := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return TimeData{Time: startUTC}
}

func (t TimeData) EndOfDay() TimeData {
	utcTime := t.Time.UTC()
	year, month, day := utcTime.Date()
	endUTC := time.Date(year, month, day, 23, 59, 59, 999999999, time.UTC)
	return TimeData{Time: endUTC}
}

func (t TimeData) TruncateHour() TimeData {
	utcTime := t.Time.UTC()
	year, month, day := utcTime.Date()
	hour := utcTime.Hour()
	truncUTC := time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
	return TimeData{Time: truncUTC}
}

func (t TimeData) StartOfMonth() TimeData {
	utcTime := t.Time.UTC()
	year, month, _ := utcTime.Date()
	startUTC := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return TimeData{Time: startUTC}
}

func TimeNow() TimeData {
	nowUTC := time.Now().UTC()
	return TimeData{Time: nowUTC}
}

func TimeNowHourly() TimeData {
	nowUTC := time.Now().UTC()
	year, month, day := nowUTC.Date()
	hour := nowUTC.Hour()
	truncUTC := time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
	return TimeData{Time: truncUTC}
}

func TimeNowDaily() TimeData {
	nowUTC := time.Now().UTC()
	year, month, day := nowUTC.Date()
	startUTC := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return TimeData{Time: startUTC}
}

func ParseDate(dateStr string) (TimeData, error) {
	if dateStr == "" {
		return TimeData{}, errors.New("date string is empty")
	}

	var parsedTime time.Time
	var err error

	parsedTime, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return TimeData{}, err
		}
	}

	return TimeData{Time: parsedTime.UTC()}, nil
}

func NewTimeData(t time.Time) TimeData {
	return TimeData{Time: t.UTC()}
}

func (t TimeData) IsFirstDayOfMonth() bool {
	dateStartOfDay := t.StartOfDay()
	dateStartOfMonth := t.StartOfMonth()
	return dateStartOfDay.Time.Equal(dateStartOfMonth.Time)
}

func (t TimeData) GetMonthlyRangeDates() (startDate TimeData, endDate TimeData) {
	utcTime := t.Time.UTC()
	year, month, _ := utcTime.Date()

	startUTC := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	startDate = TimeData{Time: startUTC}

	endDate = t.StartOfDay()
	return startDate, endDate
}

func (t TimeData) Add(d time.Duration) TimeData {
	return TimeData{Time: t.Time.UTC().Add(d)}
}

func (t TimeData) Truncate(d time.Duration) TimeData {
	return TimeData{Time: t.Time.UTC().Truncate(d)}
}

func (t TimeData) FormatLayout(layout string) string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Time.UTC().Format(layout)
}

func TimeUntil(t TimeData) time.Duration {
	now := time.Now().UTC()
	return t.Time.UTC().Sub(now)
}

func TimeSince(t TimeData) time.Duration {
	now := time.Now().UTC()
	return now.Sub(t.Time.UTC())
}
