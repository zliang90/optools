package ipapi

import (
	"fmt"
	"io"

	"github.com/zliang90/optools/geoip/core/common"
	"github.com/zliang90/optools/geoip/util"
)

const (
	ipApiJSON = "http://ip-api.com/json/%s"
)

type IPAPI struct {
	*common.Geo
}

func NewIPAPI(w io.Writer, ips []string, o *common.Option) *IPAPI {
	i := &IPAPI{
		Geo: common.NewGeo(getGeoIPInfo, w, ips, o),
	}

	i.SetHeader([]string{"IP", "As", "Country Code", "Country", "Region Name", "Isp", "Org", "TimeZone/Location", "Zip", "Comment"})

	return i
}

func (i *IPAPI) Display() {
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

		//i.SetHeader([]string{"IP", "As", "Country Code", "Country", "Region Name", "Isp", "Org", "TimeZone/Location", "Comment"})
		rows = append(rows, []string{
			geo.IP,
			geo.As,
			geo.CountryCode,
			geo.Country,
			geo.RegionName,
			geo.Isp,
			geo.Org,
			fmt.Sprintf("%s(%v,%v)", geo.TimeZone, geo.Lat, geo.Lon),
			geo.Zip,
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
	r.IP = ip

	return &r, nil
}
