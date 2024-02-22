package taobao

import (
	"fmt"
	"io"

	"github.com/zliang90/optools/geoip/core"
	"github.com/zliang90/optools/geoip/core/common"
	"github.com/zliang90/optools/geoip/util"
)

const (
	ipApiJSON = "https://ip.taobao.com/outGetIpInfo?ip=%s&accessKey=alibaba-inc"
)

type IPTaoBao struct {
	*common.Geo
}

var _ core.Interface = &IPTaoBao{}

func NewIPTaoBao(w io.Writer, ips []string, o *common.Option) *IPTaoBao {
	i := &IPTaoBao{
		Geo: common.NewGeo(getGeoIPInfo, w, ips, o),
	}

	i.SetHeader([]string{"IP", "Country", "Country Id", "Region", "Region Id", "City", "City Id", "ISP", "Comment"})

	return i
}

func (i *IPTaoBao) Display() {
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
			geo.Country,
			geo.CountryId,
			geo.Region,
			geo.RegionId,
			geo.City,
			geo.CityId,
			geo.ISP,
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
	if (r.Code != nil && *r.Code != 0) && (r.Msg != nil && *r.Msg != "") {
		return nil, fmt.Errorf("unexpected code: %d, msg: %q", *r.Code, *r.Msg)
	}
	if r.Data == nil {
		return nil, fmt.Errorf("no iptaobao geo data found for ip: %q", ip)
	}
	return r.Data, nil
}
