package lite2mmdb

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/zliang90/optools/geoip/core"
	"github.com/zliang90/optools/geoip/core/common"
)

/*

GeoLite.mmdb:
https://github.com/P3TERX/GeoLite.mmdb

Google maps:
https://www.google.com/maps
https://developers.google.com/maps/documentation/geocoding/requests-reverse-geocoding?hl=zh-cn

download sources:
	"https://git.io"
	"https://github.com/P3TERX/GeoLite.mmdb/raw/download"


*/

const (
	asn     = "asn"
	country = "country"
	city    = "city"
)

var lite2DBFiles = map[string]string{
	//asn:     "GeoLite2-ASN.mmdb",
	//country: "GeoLite2-Country.mmdb",
	city: "GeoLite2-City.mmdb",
}

//go:embed *.mmdb
var dbFs embed.FS

var dbReaders = make(map[string]*geoip2.Reader)

type Lite2MMDB struct {
	*common.Geo
}

var _ core.Interface = &Lite2MMDB{}

func NewLite2MMDB(w io.Writer, ips []string, o *common.Option) *Lite2MMDB {
	i := &Lite2MMDB{
		Geo: common.NewGeo(getGeoIPInfo, w, ips, o),
	}

	i.SetHeader([]string{"IP", "Country Code", "Country", "Continent", "City", "TimeZone/Location", "Zip", "Comment"})

	return i
}

func (l *Lite2MMDB) Init() error {
	// opendbs
	for alias, dbFileName := range lite2DBFiles {
		// read geo db file from embed.FS
		Bytes, err := dbFs.ReadFile(dbFileName)
		if err != nil {
			return err
		}
		// get geo db reader object from bytes
		r, err := geoip2.FromBytes(Bytes)
		if err != nil {
			return err
		}
		dbReaders[alias] = r
	}

	return nil
}

func (l *Lite2MMDB) Close() error {
	if dbReaders == nil || len(dbReaders) == 0 {
		return errors.New("lite2mmdbs does not opened")
	}

	for _, db := range dbReaders {
		if db != nil {
			if err := db.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *Lite2MMDB) Display() {
	if l.OutPrintJSON() {
		return
	}

	rows := make([][]string, 0)

	for _, geoIp := range l.GeoIPs {
		var comment string

		geo := geoIp.(*GeoIP)
		//fmt.Printf("geo: %v\n", util.PrettyJSON(geo))

		if geo.IP == l.MyExternalAddr {
			comment = "My External IP"
		}
		rows = append(rows, []string{
			geo.IP,
			geo.Country.IsoCode,
			geo.Country.Names["en"],
			geo.City.Continent.Names["en"],
			geo.City.City.Names["en"],
			fmt.Sprintf("%s(%v,%v)", geo.Location.TimeZone, geo.Location.Latitude, geo.Location.Longitude),
			geo.Postal.Code,
			comment,
		})
	}

	l.SetDataRows(rows)
	l.TableWrite.Display()
}

func getGeoIPInfo(ip string) (any, error) {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, fmt.Errorf("illegal ip: %q, can not to parse IP", ip)
	}

	// returns geo ip
	geo := GeoIP{
		IP: ip,
	}

	if db, ok := dbReaders[city]; ok && db != nil {
		_city, err := db.City(ipAddr)
		if err != nil {
			return nil, fmt.Errorf("parse ip error: %v", err)
		}
		geo.City = _city
	}

	return &geo, nil
}
