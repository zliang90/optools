package ipcc

/*
▶ curl -H "Accept: application/json" https://api.ip.cc/54.196.99.49/json | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   341  100   341    0     0    941      0 --:--:-- --:--:-- --:--:--   994
{
  "code": 1,
  "msg": "success",
  "ip": "54.196.99.49",
  "data": {
    "continent": "北美洲",
    "continent_code": "NA",
    "country": "美国",
    "country_code": "US",
    "region": "弗吉尼亚州",
    "region_code": "VA",
    "city": "勞登縣",
    "zip": "20149",
    "timezone": "America/New_York",
    "latitude": 39.0438,
    "longitude": -77.4874,
    "isp": "Amazon.com, Inc.",
    "user_type": "数据中心"
  }
}
*/

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	IP   string `json:"ip"`

	Data *GeoIP `json:"data"`
}

type GeoIP struct {
	IP            string  `json:"ip"`
	City          string  `json:"city"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Isp           string  `json:"isp"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	TimeZone      string  `json:"timezone"`
	UserType      string  `json:"user_type"`
	Zip           string  `json:"zip"`
}
