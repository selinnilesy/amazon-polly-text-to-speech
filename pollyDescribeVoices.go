package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"

	"fmt"
	"os"
)

func main() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create Polly client
	svc := polly.New(sess)

	// Get voices for US English
	accents := []string{"en-US", "en-GB"}
	for acc := range accents {

		input := &polly.DescribeVoicesInput{LanguageCode: aws.String(acc)}

		resp, err := svc.DescribeVoices(input)
		if err != nil {
			fmt.Println("Got error calling DescribeVoices:")
			fmt.Print(err.Error())
			os.Exit(1)
		}

		for _, v := range resp.Voices {
			fmt.Println("Name:   " + *v.Name)
			fmt.Println("Gender: " + *v.Gender)
			fmt.Println("")
		}
	}

}
