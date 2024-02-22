package ipinfo

/*

▶ curl https://ipinfo.io/8.8.8.8/json
{
  "ip": "8.8.8.8",
  "hostname": "dns.google",
  "anycast": true,
  "city": "Mountain View",
  "region": "California",
  "country": "US",
  "loc": "37.4056,-122.0775",
  "org": "AS15169 Google LLC",
  "postal": "94043",
  "timezone": "America/Los_Angeles",
  "readme": "https://ipinfo.io/missingauth"
}

*/

type GeoIP struct {
	IP       string `json:"ip"`
	Anycast  bool   `json:"anycast"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Hostname string `json:"hostname"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Readme   string `json:"readme"`
	Region   string `json:"region"`
	TimeZone string `json:"timezone"`
}
