package ip2location

import (
	"fmt"
	"io"
	"strconv"

	"github.com/zliang90/optools/geoip/core"
	"github.com/zliang90/optools/geoip/core/common"
	"github.com/zliang90/optools/geoip/util"
)

const (
	ipApiJSON = "https://api.ip2location.io/?ip=%s"
)

type IP2Location struct {
	*common.Geo
}

var _ core.Interface = &IP2Location{}

func NewIP2Location(w io.Writer, ips []string, o *common.Option) *IP2Location {
	i := &IP2Location{
		Geo: common.NewGeo(getGeoIPInfo, w, ips, o),
	}

	i.SetHeader([]string{"IP", "Asn", "As", "Country Name", "Country Code", "Region Name", "City Name", "TimeZone/Location", "Zip", "Comment"})

	return i
}

func (i *IP2Location) Display() {
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

		rows = append(rows, []string{
			geo.IP,
			strconv.Itoa(geo.Asn),
			geo.As,
			geo.CountryName,
			geo.CountryCode,
			geo.RegionName,
			geo.CityName,
			fmt.Sprintf("%s(%v,%v)", geo.TimeZone, geo.Latitude, geo.Longitude),
			geo.ZipCode,
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

	// FixMe: Hong Kong belongs to China!!!
	if r.CountryCode == "HK" {
		r.CountryName = "China"
	}

	return &r, nil
}
