package kit

import "time"

const (
	DstSwitchModeMarchSecondSundayNovemberFirstSunday         = "m1"
	DstSwitchModeMarchLastSundayOctoberLastSunday             = "m2"
	DstSwitchModeOctoberFirstSundayAprilFirstSunday           = "m3"
	DstSwitchModeSeptemberFirstSundayAprilFirstSunday         = "m4"
	DstSwitchModeMarchFridayBeforeLastSundayOctoberLastSunday = "m5"
	DstSwitchMode22March21September                           = "m6"
)

type Country struct {
	NameEng        string      // NameEng country name in english
	TranslationKey string      // TranslationKey country translation key
	Code           string      // Code country code
	Alfa2          string      // Alfa2 country alfa-2 code
	Alfa3          string      // Alfa3 country alfa-3 code
	TimeZones      []string    // TimeZones related time zones
	PhoneCodes     []string    // PhoneCodes related phone codes
	Currencies     []*Currency // Currencies related currencies
	DstSwitchMode  string      // DstSwitchMode specifies how DST (daylight saving time support) is enabled based on predefined modes. if empty, DST isn't supported
}

var (
	countriesByAlfa2 = map[string]*Country{}
	countriesByAlfa3 = map[string]*Country{}
	countriesByCode  = map[string]*Country{}
)

func init() {
	// load countries maps
	for _, c := range countries {
		countriesByCode[c.Code] = c
		countriesByAlfa2[c.Alfa2] = c
		countriesByAlfa3[c.Alfa3] = c
	}
	// load timezones
	for k, v := range tzOffsets {
		tzLocations[k] = time.FixedZone(k, int(v.Seconds()))
	}
}

// GetCountryByAlfa2 retrieves country by Alfa2 code
func GetCountryByAlfa2(alfa2 string) *Country {
	return countriesByAlfa2[alfa2]
}

// Alfa2Valid checks if alfa-2 code is valid and supported
func Alfa2Valid(alfa2 string) bool {
	return GetCountryByAlfa2(alfa2) != nil
}

// GetCountryByAlfa3 retrieves country by Alfa3 code
func GetCountryByAlfa3(alfa3 string) *Country {
	return countriesByAlfa3[alfa3]
}

func Alfa3ToAlfa2Code(alfa3 string) string {
	if alfa3 == "" {
		return ""
	}
	country := GetCountryByAlfa3(alfa3)
	if country == nil {
		return ""
	}
	return country.Alfa2
}

func GetCountryNameByAlfa3(alfa3 string) string {
	if alfa3 == "" {
		return ""
	}
	country := GetCountryByAlfa3(alfa3)
	if country == nil {
		return ""
	}
	return country.NameEng
}

func GetFirstCountryCurrencyByAlfa2(alfa2 string) string {
	if alfa2 == "" {
		return ""
	}
	country := GetCountryByAlfa2(alfa2)
	if country == nil {
		return ""
	}
	return country.Currencies[0].IsoCode
}

func Alfa2ToAlfa3Code(alfa2 string) string {
	if alfa2 == "" {
		return ""
	}
	country := GetCountryByAlfa2(alfa2)
	if country == nil {
		return ""
	}
	return country.Alfa3
}

// Alfa3Valid checks if alfa-3 code is valid and supported
func Alfa3Valid(alfa3 string) bool {
	return GetCountryByAlfa3(alfa3) != nil
}

// GetCountryByCode retrieves country by code
func GetCountryByCode(code string) *Country {
	return countriesByCode[code]
}

func CountryCodeValid(code string) bool {
	return GetCountryByCode(code) != nil
}

// GetCountries retrieves all countries
func GetCountries() []*Country {
	return countries
}

// TzValid returns true if time-zone is valid for the country
func (c *Country) TzValid(tz string) bool {
	for _, t := range c.TimeZones {
		if tz == t {
			return true
		}
	}
	return false
}

// CurrencyValid returns true if currency is valid for the country
func (c *Country) CurrencyValid(curCode string) bool {
	for _, cur := range c.Currencies {
		if cur.IsoCode == curCode {
			return true
		}
	}
	return false
}
