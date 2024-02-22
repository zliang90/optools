package ipcc

import (
	"fmt"
	"io"

	"github.com/zliang90/optools/geoip/core"
	"github.com/zliang90/optools/geoip/core/common"
	"github.com/zliang90/optools/geoip/util"
)

const (
	ipApiJSON = "https://api.ip.cc/%s/json"
)

type IPCC struct {
	*common.Geo
}

var _ core.Interface = &IPCC{}

func NewIPCC(w io.Writer, ips []string, o *common.Option) *IPCC {
	i := &IPCC{
		Geo: common.NewGeo(getGeoIPInfo, w, ips, o),
	}

	i.SetHeader([]string{"IP", "Continent", "Country", "Region", "City", "Isp", "Time Zone/Location", "Zip", "User Type", "Comment"})

	return i
}

func (i *IPCC) Display() {
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
			geo.Continent,
			geo.Country,
			geo.Region,
			geo.City,
			geo.Isp,
			fmt.Sprintf("%s(%v,%v)", geo.TimeZone, geo.Latitude, geo.Longitude),
			geo.Zip,
			geo.UserType,
			comment,
		})
	}

	i.SetDataRows(rows)
	i.TableWrite.Display()
}

func getGeoIPInfo(ip string) (any, error) {
	url := fmt.Sprintf(ipApiJSON, ip)

	var r Response
	_, err := util.RequestJSON("GET", url, nil, nil, &r)
	if err != nil {
		return nil, err
	}
	if r.Code != 1 && r.Msg != "" {
		return nil, fmt.Errorf("unexpected api response for ip %s, code: %d, msg: %s",
			ip, r.Code, r.Msg)
	}
	if r.Data == nil {
		return nil, fmt.Errorf("no ipcc geo data found for ip: %q", ip)
	}
	r.Data.IP = ip

	return r.Data, nil
}
