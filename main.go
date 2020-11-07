package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type dnsRecord struct {
	Type   string
	Values []string
}

func main() {
	s := session.Must(session.NewSession())
	r := route53.New(s)

	ids, err := getPublicZoneIds(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range ids {
		getHostedZoneData(r, v)
	}
}

func getPublicZoneIds(r *route53.Route53) ([]string, error) {
	var zones []string
	input := &route53.ListHostedZonesInput{}

	for {
		res, err := r.ListHostedZones(input)
		if err != nil {
			return nil, err
		}

		for _, v := range res.HostedZones {
			i := *v.Id
			// only audit public zones
			if !*v.Config.PrivateZone {
				// assume zone IDs start with "Z"
				if !strings.HasPrefix(i, "Z") {
					log.Printf("Skipping malformed zone ID: %v", i)
					continue
				}
				id := strings.Split(i, "/")[2]
				zones = append(zones, id)
			}
		}

		if *res.IsTruncated {
			input = &route53.ListHostedZonesInput{
				Marker: aws.String(*res.NextMarker),
			}
		} else {
			return zones, nil
		}
	}
}

func getHostedZoneData(r *route53.Route53, id string) map[string]dnsRecord {
	var result map[string]dnsRecord

	input := &route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(id)}
	res, err := r.ListResourceRecordSets(input)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v", res)
	for _, v := range res.ResourceRecordSets {
		fmt.Printf("%v,%v,%v\n", *v.Name, *v.Type, v.ResourceRecords)
	}
}
