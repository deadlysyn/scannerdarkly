package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/spf13/viper"
)

type dnsRecord struct {
	Name   string
	Type   string
	Alias  bool
	Values []string
	Active []string
}

var (
	DB = make(map[string][]dnsRecord)
)

func populateDB(ctx context.Context, r *route53.Client, zoneIDs []string) {

	for _, v := range zoneIDs {
		getResourceRecords(ctx, r, v)
	}

	count := 0
	for k := range DB {
		count = count + len(DB[k])
	}
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

func scan() {
	for id, recs := range DB {
		for idx := range recs {
			scanTCP(&DB[id][idx])
		}
	}
}

func scanTCP(rec *dnsRecord) {
	ports := viper.GetStringSlice("ports")
	timeout := (time.Duration(viper.GetInt("timeout")) * time.Second)

	for _, v := range rec.Values {
		for _, p := range ports {
			host := fmt.Sprintf("%v:%v", v, p)
			if rec.Type == "AAAA" && !rec.Alias {
				host = fmt.Sprintf("[%v]:%v", v, p)
			}
			fmt.Printf("Scanning %v...", host)
			conn, err := net.DialTimeout("tcp", host, timeout)
			if err != nil {
				fmt.Println(" closed.")
				continue
			}
			defer conn.Close()
			fmt.Println(" open.")
			rec.Active = append(rec.Active, host)
			break
		}
	}
}

func reportCSV() {
	fmt.Printf("Writing report: %v\n", name)

	f, err := os.Create(name)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// header
	w.Write([]string{
		"Zone ID",
		"Name",
		"Type",
		"Values (no open ports)",
	})

	for id, recs := range DB {
		for _, rec := range recs {
			if len(rec.Active) != 0 {
				continue
			}

			t := rec.Type
			if rec.Alias {
				t = "Alias"
			}

			w.Write([]string{
				id,
				rec.Name,
				t,
				rec.Values[0],
			})

			for i := 1; i < len(rec.Values); i++ {
				w.Write([]string{
					"",
					"",
					"",
					rec.Values[i],
				})
			}
		}
	}
}

func reportJSON() {
	json, err := json.MarshalIndent(DB, "", "\t")
	if err != nil {
		log.Fatalf(err.Error())

	}

	fmt.Println(string(json))
}
