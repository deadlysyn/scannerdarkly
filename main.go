package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var S *session.Session

func init() {
	S = session.Must(session.NewSession())
}

func main() {
	r := route53.New(S)

	getPublicZones(r)

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

func getPublicZones(r *route53.Route53) {
	input := &route53.ListHostedZonesInput{}
	res, err := r.ListHostedZones(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*res.IsTruncated)
}
