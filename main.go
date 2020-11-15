package main

import (
	"log"

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

func main() {
	s := session.Must(session.NewSession())
	r := route53.New(s)

	ids, err := getPublicZoneIds(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range ids {
		// populates DB
		getResourceRecords(r, v)
	}

	// fmt.Printf("%+v", DB)
	report()
}
