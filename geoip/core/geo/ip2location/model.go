package ip2location

/*
▶ curl -sS 'https://api.ip2location.io/?ip=8.8.8.8' | jq .
{
  "ip": "8.8.8.8",
  "country_code": "US",
  "country_name": "United States of America",
  "region_name": "California",
  "city_name": "Mountain View",
  "latitude": 37.38605,
  "longitude": -122.08385,
  "zip_code": "94035",
  "time_zone": "-08:00",
  "asn": 15169,
  "as": "Google LLC",
  "is_proxy": false,
  "message": "Limit to 500 queries per day. Sign up for a Free plan at https://www.ip2location.io to get 30K queries per month."
}
*/

type GeoIP struct {
	IP          string  `json:"ip"`
	As          string  `json:"as"`
	Asn         int     `json:"asn"`
	CityName    string  `json:"city_name"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	IsProxy     bool    `json:"is_proxy"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Message     string  `json:"message"`
	RegionName  string  `json:"region_name"`
	TimeZone    string  `json:"time_zone"`
	ZipCode     string  `json:"zip_code"`
}
