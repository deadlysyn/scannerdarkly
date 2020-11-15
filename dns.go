package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func getPublicZoneIds(r *route53.Route53) ([]string, error) {
	var ids []string
	input := &route53.ListHostedZonesInput{}

	for {
		res, err := r.ListHostedZones(input)
		if err != nil {
			return nil, err
		}

		for _, v := range res.HostedZones {
			id := *v.Id
			// only audit public zones
			if !*v.Config.PrivateZone {
				id = strings.Split(id, "/")[2]
				// assume zone IDs should start with "Z"
				if !strings.HasPrefix(id, "Z") {
					log.Printf("Skipping malformed zone ID: %v", id)
					continue
				}
				ids = append(ids, id)
			}
		}

		if *res.IsTruncated {
			input = &route53.ListHostedZonesInput{
				Marker: aws.String(*res.NextMarker),
			}
		} else {
			return ids, nil
		}
	}
}

func getResourceRecords(r *route53.Route53, id string) {
	input := &route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(id)}
	res, err := r.ListResourceRecordSets(input)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range res.ResourceRecordSets {
		rec := dnsRecord{
			ID: id,
		}
		fmt.Printf("\tTYPE: %+v\n", *s.Type)
		rec.Type = *s.Type
		if s.AliasTarget != nil {
			fmt.Printf("\t\tVALUE: %+v\n", *s.AliasTarget.DNSName)
			rec.Values = append(rec.Values, *s.AliasTarget.DNSName)
			continue
		}
		for _, r := range s.ResourceRecords {
			fmt.Printf("\t\tVALUE: %+v\n", *r.Value)
			rec.Values = append(rec.Values, *r.Value)
		}
		DB[*s.Name] = rec
	}
}
