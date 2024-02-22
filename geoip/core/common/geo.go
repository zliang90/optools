package common

import (
	"io"

	"github.com/zliang90/optools/geoip/util"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
)

type IPHandler func(ip string) (any, error)

type Geo struct {
	GeoIPs []any

	iphandle IPHandler

	*Option

	*TableWrite
}

func NewGeo(iphandle IPHandler, w io.Writer, ips []string, o *Option) *Geo {
	if o == nil {
		o = DefaultOption()
	}

	g := &Geo{
		iphandle:   iphandle,
		Option:     o,
		TableWrite: NewTableWrite(w, o),
	}

	if ips != nil {
		g.IpSets = sets.NewString(ips...)
	}

	return g
}

func (*Geo) Init() error {
	return nil
}

func (*Geo) Close() error {
	return nil
}

func (g *Geo) Load() error {
	if g.ShowMyExternal {
		ip := util.MyExternalIPAddr()
		geo, err := g.iphandle(ip)
		if err != nil {
			klog.Warningf("can not load %q geo data: %v", ip, err)
		}

		if geo != nil {
			g.MyExternalAddr = ip
			g.GeoIPs = append(g.GeoIPs, geo)
		}
		g.IpSets.Delete(ip)
	}

	for _, ip := range g.IpSets.List() {
		geo, err := g.iphandle(ip)
		if err != nil {
			klog.Warningf("can not load %q geo data: %v", ip, err)
			continue
		}
		g.GeoIPs = append(g.GeoIPs, geo)
	}

	return nil
}

func (g *Geo) OutPrintJSON() bool {
	if g.OutToJSON {
		println(util.PrettyJSON(g.GeoIPs))
		return true
	}
	return false
}
