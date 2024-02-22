package taobao

/*
▶ curl -sS 'https://ip.taobao.com/outGetIpInfo?ip=202.106.0.20&accessKey=alibaba-inc' |jq .
{
  "data": {
    "area": "",
    "country": "中国",
    "isp_id": "100026",
    "queryIp": "202.106.0.20",
    "city": "北京",
    "ip": "202.106.0.20",
    "isp": "联通",
    "county": "",
    "region_id": "110000",
    "area_id": "",
    "county_id": null,
    "region": "北京",
    "country_id": "CN",
    "city_id": "110100"
  },
  "msg": "query success",
  "code": 0
}
*/
type Response struct {
	Code *int    `json:"code"`
	Msg  *string `json:"msg"`

	Data *GeoIP `json:"data,omitempty"`
}

type GeoIP struct {
	IP        string `json:"ip"`
	Area      string `json:"area"`
	RegionId  string `json:"region_id"`
	Region    string `json:"region"`
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	ISP       string `json:"isp"`
	ISPId     string `json:"isp_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
}
