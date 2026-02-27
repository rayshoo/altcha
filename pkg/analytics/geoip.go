package analytics

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type GeoIP struct {
	db *maxminddb.Reader
}

type geoRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	Continent struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"continent"`
}

func NewGeoIP(path string) (*GeoIP, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return nil, err
	}
	return &GeoIP{db: db}, nil
}

func (g *GeoIP) Lookup(ipStr string) (country *string, continent *string) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, nil
	}

	var rec geoRecord
	if err := g.db.Lookup(ip, &rec); err != nil {
		return nil, nil
	}

	if rec.Country.ISOCode != "" {
		country = &rec.Country.ISOCode
	}
	if rec.Continent.Code != "" {
		continent = &rec.Continent.Code
	}
	return
}

func (g *GeoIP) Close() {
	g.db.Close()
}
