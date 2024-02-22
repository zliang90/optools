package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zliang90/optools/geoip/core"
	"github.com/zliang90/optools/geoip/core/common"
	"github.com/zliang90/optools/geoip/core/geo/ip2location"
	"github.com/zliang90/optools/geoip/core/geo/ipapi"
	"github.com/zliang90/optools/geoip/core/geo/ipcc"
	"github.com/zliang90/optools/geoip/core/geo/ipinfo"
	"github.com/zliang90/optools/geoip/core/geo/lite2mmdb"
	"github.com/zliang90/optools/geoip/core/geo/taobao"
	"k8s.io/klog/v2"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	GeoLite2MMDB = "lite2mmdb"

	GeoIPcc = "ipcc"

	GeoIPTaoBao = "iptaobao"

	GeoIP2Location = "ip2location"

	GeoIPAPI = "ipapi"

	GeoIPInfo = "ipinfo"
)

var (
	geoType string

	disableMyExternal     bool
	disableAutoMergeCells bool
	rowLine               bool
	alignCenter           bool
	outToJSON             bool
)

func init() {
	klog.InitFlags(nil)
}

func loadGeo(geoType string, w io.Writer, ips []string) (core.Interface, error) {
	o := &common.Option{
		ShowMyExternal: !disableMyExternal,
		OutToJSON:      outToJSON,

		TableRowLine:        rowLine,
		TableAlignCenter:    alignCenter,
		TableAutoMergeCells: !disableAutoMergeCells,
		IpSets:              sets.NewString(ips...),
	}

	var i core.Interface

	switch geoType {
	case GeoLite2MMDB:
		i = lite2mmdb.NewLite2MMDB(w, ips, o)
	case GeoIPcc:
		// return nil, fmt.Errorf("%q not implemented", geoType)
		i = ipcc.NewIPCC(w, ips, o)
	case GeoIPTaoBao:
		// o.ShowMyExternal = false
		i = taobao.NewIPTaoBao(w, ips, o)
	case GeoIP2Location:
		i = ip2location.NewIP2Location(w, ips, o)
	case GeoIPAPI:
		i = ipapi.NewIPAPI(w, ips, o)
	case GeoIPInfo:
		i = ipinfo.NewIPInfo(w, ips, o)
	default:
		return nil, fmt.Errorf("unsupported geo type: %q", geoType)
	}

	if err := i.Init(); err != nil {
		return nil, err
	}

	return i, nil
}

func addFlags() {
	pflag.StringVarP(&geoType, "geo", "g", GeoLite2MMDB, "use sourcees for ip geo, eg: lite2mmdb,ipcc,iptaobao,ipapi,ipinfo or ip2location")
	pflag.BoolVar(&disableMyExternal, "disable-my-external", false, "disable show my external ip address")
	pflag.BoolVar(&disableAutoMergeCells, "disable-auto-merge-cells", false, "disable auto merge cells")
	pflag.BoolVar(&rowLine, "row-line", true, "show row line")
	pflag.BoolVar(&alignCenter, "set-align-center", false, "set column align center")
	pflag.BoolVar(&outToJSON, "json", false, "output with json format, default table")
}

func main() {
	addFlags()
	pflag.Parse()

	i, err := loadGeo(geoType, os.Stdout, pflag.Args())
	if err != nil {
		klog.Error(err)
		return
	}
	defer i.Close()

	if err = i.Load(); err != nil {
		klog.Error(err)
		return
	}

	i.Display()
}
