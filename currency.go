package kit

type Currency struct {
	NameEng        string // NameEng currency name in english
	TranslationKey string // TranslationKey currency translation key
	IsoCode        string // IsoCode currency ISO code
	Number         string // Number currency number code
	Symbol         string // Symbol currency symbol
	Unit           int    // Unit used mostly in UI operate with different currencies regardless its value (aprox it equal 1 eur)
}

// GetCurrencyCodes retrieves currency by code
func GetCurrencyCodes() []string {
	return MapSet(currenciesByISO, func(key string, item *Currency) string {
		return key
	})
}

// GetCurrency retrieves currency by code
func GetCurrency(code string) *Currency {
	return currenciesByISO[code]
}

// GetCurrencies retrieves currencies by codes
func GetCurrencies(codes ...string) []*Currency {
	var r []*Currency
	for _, c := range codes {
		cur := currenciesByISO[c]
		if cur != nil {
			r = append(r, cur)
		}
	}
	return r
}

// CurrencyValid checks if currency valid and supported
func CurrencyValid(code string) bool {
	return currenciesByISO[code] != nil
}
