package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type dnsRecord struct {
	Name   string
	Type   string
	Values []string
}

// DB maps Route53 Zone IDs to slices of dnsRecords
var DB = make(map[string][]dnsRecord)

// PORTS ports to check on DNS targets (default: 80 8080 443 8443)
var PORTS []string

// ZONES Route53 zone IDs to audit (default: all in account)
var ZONES []string

func init() {
	ports := strings.TrimSpace(os.Getenv("PORTS"))
	if len(ports) > 0 {
		for _, p := range strings.Split(ports, " ") {
			PORTS = append(PORTS, p)
		}
	} else {
		PORTS = []string{
			"80",
			"8080",
			"443",
			"8443",
		}
	}

	zones := strings.TrimSpace(os.Getenv("ZONES"))
	if len(zones) > 0 {
		for _, z := range strings.Split(zones, " ") {
			ZONES = append(ZONES, z)
		}
	}
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

	// scan()
	reportCSV()
	// reportJSON()
}

func populateDB(r *route53.Route53, ids []string) {
	for _, v := range ids {
		getResourceRecords(r, v)
	}
}
