package ipinfo

import (
	"fmt"
	"io"

	"github.com/zliang90/optools/geoip/core/common"
	"github.com/zliang90/optools/geoip/util"
)

const (
	ipApiJSON = "https://ipinfo.io/%s/json"
)

type IPInfo struct {
	*common.Geo
}

func NewIPInfo(w io.Writer, ips []string, o *common.Option) *IPInfo {
	i := &IPInfo{
		Geo: common.NewGeo(getGeoIPInfo, w, ips, o),
	}

	i.SetHeader([]string{"IP", "HostName", "Country", "Region", "City", "Org", "TimeZone/Location", "Zip", "Comment"})

	return i
}

func (i *IPInfo) Display() {
	if i.OutPrintJSON() {
		return
	}

	rows := make([][]string, 0)

	for _, geoIp := range i.GeoIPs {
		var comment string

		geo := geoIp.(*GeoIP)

		if geo.IP == i.MyExternalAddr {
			comment = "My External IP"
		}

		//i.SetHeader([]string{"IP", "HostName", "Country", "Region", "City", "Org", "TimeZone/Location", "Zip", "Comment"})
		rows = append(rows, []string{
			geo.IP,
			geo.Hostname,
			geo.Country,
			geo.Region,
			geo.City,
			geo.Org,
			fmt.Sprintf("%s(%v)", geo.TimeZone, geo.Loc),
			geo.Postal,
			comment,
		})
	}

	i.SetDataRows(rows)
	i.TableWrite.Display()

}

func getGeoIPInfo(ip string) (any, error) {
	url := fmt.Sprintf(ipApiJSON, ip)

	var r GeoIP
	_, err := util.RequestJSON("GET", url, nil, nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
