package kit

const (
	CurRUB string = "RUB"
	CurUSD string = "USD"
	CurRSD string = "RSD"
	CurHRK string = "HRK"
	CurEUR string = "EUR"
	CurKZT string = "KZT"
	CurINR string = "INR"
	CurALL string = "ALL"
	CurUZS string = "UZS"
	CurKGS string = "KGS"
	CurAMD string = "AMD"
	CurBYN string = "BYN"
	CurCNY string = "CNY"
	CurTRY string = "TRY"
	CurBAM string = "BAM"
	CurCLP string = "CLP"
	CurGBP string = "GBP"
)

var (
	currenciesByISO = map[string]*Currency{
		CurUZS: {
			NameEng:        "Uzbekistan sum",
			TranslationKey: "currencies.uzs",
			IsoCode:        CurUZS,
			Number:         "860",
			Symbol:         "soum",
			Unit:           100,
		},
		CurKGS: {
			NameEng:        "Kyrgyzstani som",
			TranslationKey: "currencies.kgs",
			IsoCode:        CurKGS,
			Number:         "417",
			Symbol:         "som",
			Unit:           100,
		},
		CurAMD: {
			NameEng:        "Armenian dram",
			TranslationKey: "currencies.amd",
			IsoCode:        CurAMD,
			Number:         "051",
			Symbol:         "֏",
			Unit:           100,
		},
		CurBYN: {
			NameEng:        "Belarusian ruble",
			TranslationKey: "currencies.byn",
			IsoCode:        CurBYN,
			Number:         "933",
			Symbol:         "Br",
			Unit:           100,
		},
		CurCNY: {
			NameEng:        "Renminbi",
			TranslationKey: "currencies.cny",
			IsoCode:        CurCNY,
			Number:         "156",
			Symbol:         "¥",
			Unit:           10,
		},
		CurTRY: {
			NameEng:        "Turkish lira",
			TranslationKey: "currencies.try",
			IsoCode:        CurTRY,
			Number:         "949",
			Symbol:         "₺",
			Unit:           100,
		},
		CurBAM: {
			NameEng:        "Bosnia and Herzegovina convertible mark",
			TranslationKey: "currencies.bam",
			IsoCode:        CurBAM,
			Number:         "977",
			Symbol:         "KM",
			Unit:           100,
		},
		CurRUB: {
			NameEng:        "Ruble",
			TranslationKey: "currencies.rub",
			IsoCode:        CurRUB,
			Number:         "643",
			Symbol:         "₽",
			Unit:           100,
		},
		CurINR: {
			NameEng:        "Indian Rupee",
			TranslationKey: "currencies.inr",
			IsoCode:        CurINR,
			Number:         "356",
			Symbol:         "₹",
			Unit:           100,
		},
		CurALL: {
			NameEng:        "Albanian lek",
			TranslationKey: "currencies.all",
			IsoCode:        CurALL,
			Number:         "008",
			Symbol:         "Lek",
			Unit:           100,
		},
		CurHRK: {
			NameEng:        "Kuna",
			TranslationKey: "currencies.hrk",
			IsoCode:        CurHRK,
			Number:         "191",
			Symbol:         "kn",
			Unit:           10,
		},
		CurRSD: {
			NameEng:        "Dinar",
			TranslationKey: "currencies.rsd",
			IsoCode:        CurRSD,
			Number:         "941",
			Symbol:         "Дин.",
			Unit:           100,
		},
		CurUSD: {
			NameEng:        "US Dollar",
			TranslationKey: "currencies.usd",
			IsoCode:        CurUSD,
			Number:         "840",
			Symbol:         "$",
			Unit:           1,
		},
		CurEUR: {
			NameEng:        "Euro",
			TranslationKey: "currencies.eur",
			IsoCode:        CurEUR,
			Number:         "978",
			Symbol:         "€",
			Unit:           1,
		},
		CurKZT: {
			NameEng:        "Tenge",
			TranslationKey: "currencies.kzt",
			IsoCode:        CurKZT,
			Number:         "398",
			Symbol:         "₸",
			Unit:           500,
		},
		CurCLP: {
			NameEng:        "Peso",
			TranslationKey: "currencies.clp",
			IsoCode:        CurCLP,
			Number:         "152",
			Symbol:         "$",
			Unit:           1000,
		},
		CurGBP: {
			NameEng:        "Pound Sterling",
			TranslationKey: "currencies.gbp",
			IsoCode:        CurGBP,
			Number:         "826",
			Symbol:         "£",
			Unit:           1,
		},
	}
)
