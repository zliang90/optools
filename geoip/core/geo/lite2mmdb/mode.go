package lite2mmdb

import "github.com/oschwald/geoip2-golang"

type GeoIP struct {
	IP string
	*geoip2.City
}
