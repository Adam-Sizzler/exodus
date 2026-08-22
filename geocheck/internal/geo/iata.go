package geo

// iataToCountry maps airport/metro codes used by CDN point-of-presence names
// to ISO 3166-1 alpha-2 countries. It covers the cities Cloudflare, Google
// Global Cache and Netflix Open Connect actually deploy in, so the common case
// resolves offline; anything missing falls back to a network lookup.
var iataToCountry = map[string]string{
	// Europe
	"AMS": "NL", "RTM": "NL", "EIN": "NL",
	"FRA": "DE", "MUC": "DE", "DUS": "DE", "HAM": "DE", "BER": "DE", "TXL": "DE", "STR": "DE",
	"LHR": "GB", "LON": "GB", "LGW": "GB", "MAN": "GB", "EDI": "GB", "GLA": "GB", "BHX": "GB",
	"CDG": "FR", "PAR": "FR", "ORY": "FR", "MRS": "FR", "LYS": "FR", "BOD": "FR", "LIL": "FR",
	"MAD": "ES", "BCN": "ES", "VLC": "ES", "AGP": "ES", "PMI": "ES", "BIO": "ES",
	"MXP": "IT", "MIL": "IT", "LIN": "IT", "FCO": "IT", "ROM": "IT", "NAP": "IT", "PMO": "IT", "VCE": "IT", "BLQ": "IT",
	"VIE": "AT", "ZRH": "CH", "GVA": "CH", "BSL": "CH",
	"PRG": "CZ", "WAW": "PL", "KRK": "PL", "GDN": "PL", "POZ": "PL", "WRO": "PL",
	"BUD": "HU", "OTP": "RO", "BUH": "RO", "CLJ": "RO", "SOF": "BG", "ATH": "GR", "SKG": "GR",
	"IST": "TR", "SAW": "TR", "ADB": "TR", "ESB": "TR", "AYT": "TR",
	"ARN": "SE", "STO": "SE", "GOT": "SE", "MMA": "SE",
	"CPH": "DK", "OSL": "NO", "SVG": "NO", "HEL": "FI", "KEF": "IS",
	"DUB": "IE", "ORK": "IE", "LIS": "PT", "OPO": "PT", "BRU": "BE", "LUX": "LU",
	"ZAG": "HR", "BEG": "RS", "SKP": "MK", "TIA": "AL", "LJU": "SI", "BTS": "SK", "SJJ": "BA", "TGD": "ME", "PRN": "XK",
	"RIX": "LV", "VNO": "LT", "TLL": "EE",
	"SVO": "RU", "DME": "RU", "VKO": "RU", "MOW": "RU", "LED": "RU", "KZN": "RU", "SVX": "RU",
	"OVB": "RU", "KJA": "RU", "ROV": "RU", "AER": "RU", "KHV": "RU", "VVO": "RU",
	"KBP": "UA", "IEV": "UA", "ODS": "UA", "LWO": "UA", "HRK": "UA",
	"MSQ": "BY", "KIV": "MD", "TBS": "GE", "EVN": "AM", "GYD": "AZ",
	"MLA": "MT", "LCA": "CY", "NIC": "CY",

	// North America
	"IAD": "US", "DCA": "US", "BWI": "US", "WDC": "US",
	"JFK": "US", "EWR": "US", "LGA": "US", "NYC": "US", "BUF": "US", "ROC": "US",
	"ORD": "US", "MDW": "US", "CHI": "US", "IND": "US", "CMH": "US", "CLE": "US", "DTW": "US", "MKE": "US",
	"LAX": "US", "BUR": "US", "SNA": "US", "ONT": "US", "SAN": "US", "SJC": "US", "SFO": "US", "OAK": "US",
	"SEA": "US", "PDX": "US", "BOI": "US", "SLC": "US", "DEN": "US", "ABQ": "US", "PHX": "US", "TUS": "US", "LAS": "US",
	"DFW": "US", "DAL": "US", "IAH": "US", "HOU": "US", "AUS": "US", "SAT": "US", "ELP": "US", "OKC": "US", "TUL": "US",
	"ATL": "US", "MIA": "US", "TPA": "US", "MCO": "US", "JAX": "US", "CLT": "US", "RDU": "US", "RIC": "US", "ORF": "US",
	"BOS": "US", "PHL": "US", "PIT": "US", "MSP": "US", "MCI": "US", "STL": "US", "OMA": "US", "BNA": "US", "MEM": "US",
	"ANC": "US", "HNL": "US", "SJU": "PR",
	"YYZ": "CA", "YUL": "CA", "YVR": "CA", "YYC": "CA", "YEG": "CA", "YOW": "CA", "YWG": "CA", "YHZ": "CA", "YQM": "CA", "YXE": "CA",
	"MEX": "MX", "QRO": "MX", "GDL": "MX", "MTY": "MX", "TIJ": "MX", "CUN": "MX",

	// Latin America and the Caribbean
	"GRU": "BR", "SAO": "BR", "CGH": "BR", "VCP": "BR", "GIG": "BR", "RIO": "BR", "BSB": "BR",
	"POA": "BR", "CWB": "BR", "FOR": "BR", "REC": "BR", "SSA": "BR", "MAO": "BR", "CNF": "BR", "BEL": "BR", "FLN": "BR",
	"EZE": "AR", "BUE": "AR", "AEP": "AR", "COR": "AR", "SCL": "CL", "LIM": "PE",
	"BOG": "CO", "MDE": "CO", "CLO": "CO", "BAQ": "CO", "UIO": "EC", "GYE": "EC", "CCS": "VE",
	"ASU": "PY", "MVD": "UY", "LPB": "BO", "VVI": "BO",
	"PTY": "PA", "SJO": "CR", "GUA": "GT", "SAL": "SV", "TGU": "HN", "MGA": "NI", "BZE": "BZ",
	"SDQ": "DO", "STI": "DO", "KIN": "JM", "CUR": "CW", "POS": "TT", "HAV": "CU", "NAS": "BS", "BGI": "BB",

	// Asia-Pacific
	"NRT": "JP", "HND": "JP", "TYO": "JP", "KIX": "JP", "ITM": "JP", "OSA": "JP",
	"NGO": "JP", "FUK": "JP", "CTS": "JP", "OKA": "JP", "SDJ": "JP",
	"ICN": "KR", "SEL": "KR", "GMP": "KR", "PUS": "KR", "CJU": "KR",
	"HKG": "HK", "MFM": "MO", "TPE": "TW", "KHH": "TW", "TSA": "TW",
	"SIN": "SG", "KUL": "MY", "JHB": "MY", "PEN": "MY", "BKI": "MY", "KCH": "MY",
	"BKK": "TH", "DMK": "TH", "CNX": "TH", "HKT": "TH",
	"HAN": "VN", "SGN": "VN", "DAD": "VN", "MNL": "PH", "CEB": "PH", "DVO": "PH",
	"CGK": "ID", "JKT": "ID", "DPS": "ID", "SUB": "ID", "MES": "ID", "UPG": "ID",
	"PNH": "KH", "RGN": "MM", "VTE": "LA", "BWN": "BN", "DIL": "TL",
	"BOM": "IN", "DEL": "IN", "MAA": "IN", "BLR": "IN", "HYD": "IN", "CCU": "IN",
	"AMD": "IN", "PNQ": "IN", "NAG": "IN", "COK": "IN", "JAI": "IN", "LKO": "IN", "PAT": "IN", "IXC": "IN", "BBI": "IN",
	"CMB": "LK", "DAC": "BD", "CGP": "BD", "KTM": "NP", "MLE": "MV", "PBH": "BT",
	"KHI": "PK", "LHE": "PK", "ISB": "PK", "KBL": "AF",
	"PEK": "CN", "PKX": "CN", "PVG": "CN", "SHA": "CN", "CAN": "CN", "SZX": "CN", "CTU": "CN", "HGH": "CN",
	"TSN": "CN", "WUH": "CN", "XIY": "CN", "CKG": "CN", "KMG": "CN", "SHE": "CN", "TAO": "CN", "NKG": "CN",
	"FOC": "CN", "XMN": "CN", "CGO": "CN", "CSX": "CN", "HRB": "CN", "URC": "CN", "LHW": "CN",
	"ULN": "MN",
	"ALA": "KZ", "NQZ": "KZ", "TSE": "KZ", "FRU": "KG", "TAS": "UZ", "DYU": "TJ", "ASB": "TM",

	// Middle East
	"DXB": "AE", "AUH": "AE", "SHJ": "AE", "FJR": "AE",
	"DOH": "QA", "KWI": "KW", "BAH": "BH", "MCT": "OM",
	"RUH": "SA", "JED": "SA", "DMM": "SA",
	"TLV": "IL", "HFA": "IL", "AMM": "JO", "BEY": "LB", "BGW": "IQ", "EBL": "IQ", "BSR": "IQ",
	"IKA": "IR", "THR": "IR", "DAM": "SY", "SAH": "YE",

	// Africa
	"JNB": "ZA", "CPT": "ZA", "DUR": "ZA",
	"LOS": "NG", "ABV": "NG", "KAN": "NG", "ACC": "GH",
	"NBO": "KE", "MBA": "KE", "DAR": "TZ", "KGL": "RW", "EBB": "UG", "ADD": "ET", "JIB": "DJ", "MGQ": "SO",
	"CAI": "EG", "ALY": "EG", "CMN": "MA", "RBA": "MA", "TNG": "MA", "TUN": "TN", "ALG": "DZ", "TIP": "LY", "KRT": "SD",
	"DKR": "SN", "ABJ": "CI", "DLA": "CM", "NSI": "CM", "LFW": "TG", "COO": "BJ", "OUA": "BF", "BKO": "ML", "NIM": "NE",
	"LAD": "AO", "MPM": "MZ", "HRE": "ZW", "LUN": "ZM", "GBE": "BW", "WDH": "NA", "BLZ": "MW", "MRU": "MU",
	"TNR": "MG", "RUN": "RE", "SEZ": "SC", "FIH": "CD", "BZV": "CG", "LBV": "GA",

	// Oceania
	"SYD": "AU", "MEL": "AU", "BNE": "AU", "PER": "AU", "ADL": "AU", "CBR": "AU", "HBA": "AU", "DRW": "AU", "OOL": "AU",
	"AKL": "NZ", "CHC": "NZ", "WLG": "NZ",
	"NAN": "FJ", "POM": "PG", "NOU": "NC", "PPT": "PF", "GUM": "GU", "APW": "WS", "VLI": "VU",
}
