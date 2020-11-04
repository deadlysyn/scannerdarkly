package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// sess, err := session.NewSession()
	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(os.Getenv("AWS_REGION"))},
	// )
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// client instance
	// svc := s3.New(sess)
	// log calls
	// svc := dynamodb.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))

	svc := s3.New(sess)
	input := &s3.ListBucketsInput{}
	result, err := svc.ListBuckets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}
