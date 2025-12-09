package utils

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

const TimeFormat = "02-01-2006Z15:04:05Z"

type TimeData struct {
	time.Time
}

func (ct TimeData) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := ct.Time.Format(TimeFormat)
	return json.Marshal(formatted)
}

func (ct *TimeData) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	if timeStr == "" || timeStr == "null" {
		ct.Time = time.Time{}
		return nil
	}

	parsed, err := time.Parse(TimeFormat, timeStr)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return err
		}
	}

	ct.Time = parsed.UTC()
	return nil
}

func (ct TimeData) Value() (driver.Value, error) {
	if ct.Time.IsZero() {
		return nil, nil
	}
	return ct.Time, nil
}

func (ct *TimeData) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v.UTC()
		return nil
	case []byte:
		return ct.UnmarshalJSON(v)
	case string:
		return ct.UnmarshalJSON([]byte(v))
	default:
		if t, ok := value.(time.Time); ok {
			ct.Time = t.UTC()
			return nil
		}
		return nil
	}
}

func NewTimeData(t time.Time) TimeData {
	return TimeData{Time: t.UTC()}
}

func TimeNow() TimeData {
	return TimeData{Time: time.Now().UTC()}
}

func (ct TimeData) Format() string {
	if ct.Time.IsZero() {
		return ""
	}
	return ct.Time.Format(TimeFormat)
}

func TimeNowHourly() TimeData {
	t := time.Now().UTC().Truncate(time.Hour)
	return TimeData{Time: t.UTC()}
}

func TimeNowDaily() TimeData {
	now := time.Now().UTC()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return TimeData{Time: t}
}

func ParseTimeData(timeStr string) (TimeData, error) {
	if timeStr == "" {
		return TimeData{}, nil
	}

	t, err := time.Parse(TimeFormat, timeStr)
	if err != nil {
		return TimeData{}, err
	}

	return TimeData{Time: t.UTC()}, nil
}
