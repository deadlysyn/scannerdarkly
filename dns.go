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
	input := &route53.ListHostedZonesInput{
		MaxItems: aws.String("100"),
	}

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

// need to handle truncation...
func getResourceRecords(r *route53.Route53, id string) {
	var recs []dnsRecord
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(id),
		MaxItems: aws.String("100"),
	}

	for {
		res, err := r.ListResourceRecordSets(input)
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range res.ResourceRecordSets {
			rec := dnsRecord{
				Name: *s.Name,
				Type: *s.Type,
			}
			if s.AliasTarget != nil {
				rec.Values = append(rec.Values, *s.AliasTarget.DNSName)
			} else {
				for _, r := range s.ResourceRecords {
					if !strings.HasSuffix(*r.Value, "acm-validations.aws.") {
						rec.Values = append(rec.Values, *r.Value)
					} else {
						fmt.Printf("Skipping %v (ACM)\n", *s.Name)
					}
				}
			}
			if len(rec.Values) > 0 {
				recs = append(recs, rec)
			}
		}
		DB[id] = recs

		if *res.IsTruncated {
			input = &route53.ListResourceRecordSetsInput{
				StartRecordName: res.NextRecordName,
				StartRecordType: res.NextRecordType,
			}
		} else {
			break
		}
	}
}
