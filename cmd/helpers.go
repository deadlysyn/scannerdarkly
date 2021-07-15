package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func populateDB(ctx context.Context, r *route53.Client, zones []string) {
	for _, v := range zones {
		getResourceRecords(ctx, r, v)
	}
}

func getPublicZoneIds(ctx context.Context, r *route53.Client) ([]string, error) {
	var zones []string
	input := &route53.ListHostedZonesInput{
		MaxItems: aws.Int32(100),
	}

	for {
		res, err := r.ListHostedZones(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, v := range res.HostedZones {
			// only audit public zones
			if !v.Config.PrivateZone {
				id := strings.Split(*v.Id, "/")[2]
				// zone IDs should start with "Z"
				if !strings.HasPrefix(id, "Z") {
					log.Printf("Skipping malformed zone ID: %v", id)
					continue
				}
				zones = append(zones, id)
			}
		}

		if res.IsTruncated {
			input = &route53.ListHostedZonesInput{
				MaxItems: aws.Int32(100),
				Marker:   aws.String(*res.NextMarker),
			}
		} else {
			return zones, nil
		}
	}
}

func getResourceRecords(ctx context.Context, r *route53.Client, id string) {
	var recs []dnsRecord
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(id),
		MaxItems:     aws.Int32(100),
	}

	fmt.Printf("Processing zone %v\n", id)

	for {

		res, err := r.ListResourceRecordSets(ctx, input)
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range res.ResourceRecordSets {
			switch string(s.Type) {
			case "A", "AAAA", "CNAME":
				rec := dnsRecord{
					Name: strings.TrimSuffix(*s.Name, "."),
					Type: string(s.Type),
				}
				if s.AliasTarget != nil {
					rec.Alias = true
					rec.Values = append(rec.Values, strings.TrimSuffix(*s.AliasTarget.DNSName, "."))
				} else {
					for _, r := range s.ResourceRecords {
						if !strings.HasSuffix(*r.Value, "acm-validations.aws.") {
							rec.Values = append(rec.Values, strings.TrimSuffix(*r.Value, "."))
						} else {
							fmt.Fprintf(os.Stderr, "Skipping %v (ACM)\n", strings.TrimSuffix(*s.Name, "."))
						}
					}
				}
				if len(rec.Values) > 0 {
					recs = append(recs, rec)
				}
			default:
				fmt.Fprintf(os.Stderr, "Skipping %v (%v)\n", strings.TrimSuffix(*s.Name, "."), s.Type)
			}
		}
		DB[id] = recs

		if res.IsTruncated {
			input = &route53.ListResourceRecordSetsInput{
				HostedZoneId:    aws.String(id),
				MaxItems:        aws.Int32(100),
				StartRecordName: res.NextRecordName,
				StartRecordType: res.NextRecordType,
			}
		} else {
			break
		}
	}
}
