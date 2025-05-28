package ip

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
)

type ip struct {
	GeoIp2 *geoip2.Reader
}

var Ip ip

func InitIpDB(name string) {
	Ip = ip{}
	DbIp, err := geoip2.Open(name)
	Ip.GeoIp2 = DbIp
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *ip) IpIsCN(ip string) bool {
	record, err := c.GeoIp2.City(net.ParseIP(ip))
	if err != nil {
		fmt.Println(err.Error())
	}
	if record.Country.IsoCode == "CN" {
		return true
	}
	return false
}

func (c *ip) GetIpInfo(ip string) (error, *geoip2.City) {
	record, err := c.GeoIp2.City(net.ParseIP(ip))
	if err != nil {
		return err, nil
	}
	return nil, record
}
