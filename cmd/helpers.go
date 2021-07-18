package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// which record types to scan
var RRtypes = map[string]bool{
	"CNAME": true,
}

func populateDB(ctx context.Context, r *route53.Client, zoneIDs []string) {
	if scanArecords {
		RRtypes["A"] = true
		RRtypes["AAAA"] = true
	}

	for _, v := range zoneIDs {
		getResourceRecords(ctx, r, v)
	}

	count := 0
	for k := range DB {
		count = count + len(DB[k])
	}
	fmt.Printf("\n\nDEBUG: %v\n\n", count)
}

func getPublicZoneIDs(ctx context.Context, r *route53.Client) ([]string, error) {
	var zoneIDs []string

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
				ID := strings.Split(*v.Id, "/")[2]
				if !strings.HasPrefix(ID, "Z") {
					log.Printf("Skipping malformed zone ID: %v", ID)
					continue
				}
				zoneIDs = append(zoneIDs, ID)
			}
		}

		if res.IsTruncated {
			input = &route53.ListHostedZonesInput{
				MaxItems: aws.Int32(100),
				Marker:   aws.String(*res.NextMarker),
			}
		} else {
			return zoneIDs, nil
		}
	}
}

func getResourceRecords(ctx context.Context, r *route53.Client, ID string) {
	var recs []dnsRecord

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(ID),
		MaxItems:     aws.Int32(100),
	}

	fmt.Printf("Processing zone %v\n", ID)

	for {

		res, err := r.ListResourceRecordSets(ctx, input)
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range res.ResourceRecordSets {
			switch {
			case RRtypes[string(s.Type)]:

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
							fmt.Printf("\tSkipping %v (ACM)\n", strings.TrimSuffix(*s.Name, "."))
						}
					}
				}
				if len(rec.Values) > 0 {
					recs = append(recs, rec)
				}
			default:
				fmt.Printf("\tSkipping %v (%v)\n", strings.TrimSuffix(*s.Name, "."), s.Type)
			}
		}
		DB[ID] = recs

		if res.IsTruncated {
			input = &route53.ListResourceRecordSetsInput{
				HostedZoneId:    aws.String(ID),
				MaxItems:        aws.Int32(100),
				StartRecordName: res.NextRecordName,
				StartRecordType: res.NextRecordType,
			}
		} else {
			break
		}
	}
}
