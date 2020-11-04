package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	input := &route53.ListHostedZonesInput{}
	result, err := svc.ListHostedZones(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

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
