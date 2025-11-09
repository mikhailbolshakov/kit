package kit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Alfa2Valid(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{out: true, in: "RU"},
		{out: false, in: "ru"},
		{out: false, in: "rU"},
		{out: false, in: ""},
		{out: false, in: "invalid"},
		{out: false, in: "RU "},
		{out: false, in: " RU"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.out, Alfa2Valid(tt.in))
		})
	}
}

func Test_Alfa3Valid(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{out: true, in: "RUS"},
		{out: false, in: "rus"},
		{out: false, in: "rUs"},
		{out: false, in: ""},
		{out: false, in: "invalid"},
		{out: false, in: "RUS "},
		{out: false, in: " RUS"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.out, Alfa3Valid(tt.in))
		})
	}
}

func Test_CountryCodeValid(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{out: true, in: "643"},
		{out: false, in: "64 3"},
		{out: false, in: "0643"},
		{out: false, in: ""},
		{out: false, in: "invalid"},
		{out: false, in: "643 "},
		{out: false, in: " 643"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.out, CountryCodeValid(tt.in))
		})
	}
}

func Test_CountryTzValid(t *testing.T) {
	country := GetCountryByAlfa3("SRB")
	assert.True(t, country.TzValid(TzP2))
	assert.False(t, country.TzValid(TzP10))
	assert.False(t, country.TzValid(""))
	assert.False(t, country.TzValid("invalid"))
}

func Test_CountryCurrencyValid(t *testing.T) {
	country := GetCountryByAlfa3("SRB")
	assert.True(t, country.CurrencyValid(CurRSD))
	assert.False(t, country.CurrencyValid(CurRUB))
	assert.False(t, country.CurrencyValid(""))
	assert.False(t, country.CurrencyValid("invalid"))
}
