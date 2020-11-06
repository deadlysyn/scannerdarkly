package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var s *session.Session

func init() {
	s = session.Must(session.NewSession())
}

func main() {
	r := route53.New(s)

	ids, err := getPublicZoneIds(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ids)

	// input := &route53.GetHostedZoneInput{}
	// result, err := svc.GetHostedZone(input)
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		case route53.ErrCodeNoSuchHostedZone:
	// 			fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
	// 		case route53.ErrCodeInvalidInput:
	// 			fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
	// 		default:
	// 			fmt.Println(aerr.Error())
	// 		}
	// 	} else {
	// 		fmt.Println(err.Error())
	// 	}
	// 	return
	// }
	// fmt.Printf("%+v", result)
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
			id := strings.Split(*v.Id, "/")[2]
			zones = append(zones, id)
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
