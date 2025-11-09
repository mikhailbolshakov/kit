package kit

import "time"

const (
	TzUTC  = "UTC"
	TzP13  = "UTC+13"
	TzP12  = "UTC+12"
	TzP11  = "UTC+11"
	TzP10  = "UTC+10"
	TzP9   = "UTC+9"
	TzP8   = "UTC+8"
	TzP7   = "UTC+7"
	TzP6p5 = "UTC+6:30"
	TzP6   = "UTC+6"
	TzP5   = "UTC+5"
	TzP5p5 = "UTC+5:30"
	TzP4   = "UTC+4"
	TzP3   = "UTC+3"
	TzP2   = "UTC+2"
	TzP1   = "UTC+1"
	TzM1   = "UTC-1"
	TzM2   = "UTC-2"
	TzM3   = "UTC-3"
	TzM4   = "UTC-4"
	TzM5   = "UTC-5"
	TzM6   = "UTC-6"
	TzM7   = "UTC-7"
	TzM8   = "UTC-8"
	TzM9   = "UTC-9"
	TzM10  = "UTC-10"
	TzM11  = "UTC-11"
)

var (
	tzOffsets = map[string]time.Duration{
		TzUTC:  0,
		TzP13:  13 * time.Hour,
		TzP12:  12 * time.Hour,
		TzP11:  11 * time.Hour,
		TzP10:  10 * time.Hour,
		TzP9:   9 * time.Hour,
		TzP8:   8 * time.Hour,
		TzP7:   7 * time.Hour,
		TzP6p5: 6*time.Hour + 30*time.Minute,
		TzP6:   6 * time.Hour,
		TzP5:   5 * time.Hour,
		TzP5p5: 5*time.Hour + 30*time.Minute,
		TzP4:   4 * time.Hour,
		TzP3:   3 * time.Hour,
		TzP2:   2 * time.Hour,
		TzP1:   time.Hour,
		TzM1:   -time.Hour,
		TzM2:   -2 * time.Hour,
		TzM3:   -3 * time.Hour,
		TzM4:   -4 * time.Hour,
		TzM5:   -5 * time.Hour,
		TzM6:   -6 * time.Hour,
		TzM7:   -7 * time.Hour,
		TzM8:   -8 * time.Hour,
		TzM9:   -9 * time.Hour,
		TzM10:  -10 * time.Hour,
		TzM11:  -11 * time.Hour,
	}
	tzLocations = map[string]*time.Location{}
)

func TzValid(tz string) bool {
	return tzLocations[tz] != nil
}

func GetTzLocation(tz string) *time.Location {
	return tzLocations[tz]
}

func init() {
	for k, v := range tzOffsets {
		tzLocations[k] = time.FixedZone(k, int(v.Seconds()))
	}
}

func ToTz(t time.Time, tz string) (time.Time, error) {
	if tz == "" {
		return t, nil
	}
	// first check predefined locations
	loc := GetTzLocation(tz)
	if loc == nil {
		// try to load timezone
		var err error
		loc, err = time.LoadLocation(tz)
		if err != nil {
			return t, ErrNotSupportedTz(err, tz)
		}
	}
	//set timezone,
	return t.In(loc), nil
}
