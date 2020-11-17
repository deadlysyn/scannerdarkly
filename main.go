package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type dnsRecord struct {
	Name   string
	Type   string
	Alias  bool
	Values []string
	Active []string
}

// DB maps Route53 Zone IDs to slices of dnsRecords
var DB = make(map[string][]dnsRecord)

func init() {
	parseEnv()
}

func main() {
	s := session.Must(session.NewSession())
	r := route53.New(s)

	if len(ZONES) == 0 {
		ids, err := getPublicZoneIds(r)
		if err != nil {
			log.Fatal(err)
		}
		populateDB(r, ids)
	} else {
		populateDB(r, ZONES)
	}

	scan()
	reportCSV()
	// reportJSON()
}

func populateDB(r *route53.Route53, ids []string) {
	for _, v := range ids {
		getResourceRecords(r, v)
	}
}
