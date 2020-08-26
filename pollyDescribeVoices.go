package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"

	"fmt"
	"os"
)

func SynthesizeSpeechInput(text string, writeTo string) error {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create Polly client
	svc := polly.New(sess)

	//get all robot voices
	input := &polly.DescribeVoicesInput{}
	resp, err := svc.DescribeVoices(input)
	if err != nil {
		fmt.Println("Got error calling DescribeVoices:", err)
		return err
	}

	//process label with robot voices
	for _, v := range resp.Voices {
		dontWant := "en-IN"
		want := "en"
		if strings.HasPrefix(*v.LanguageCode, want) != true {
			continue
		}
		if *v.LanguageCode == dontWant {
			continue
		}

		fmt.Println("Name:   ", *v.Name)
		fmt.Println("Country:   ", *v.LanguageCode)

		var input *polly.SynthesizeSpeechInput
		if *(v.Name) == "Kevin" {
			neuralEngine := "neural"
			input = &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(text), VoiceId: v.Id, Engine: &neuralEngine}
		} else {
			input = &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(text), VoiceId: v.Id}
		}
		output, err := svc.SynthesizeSpeech(input)
		if err != nil {
			fmt.Println("Got error calling SynthesizeSpeech:", v, err.Error())
			return err
		}

		// Save as MP3
		mp3File := text + "-" + (*v.Id) + ".mp3"
		data, err := ioutil.ReadAll(output.AudioStream)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(writeTo, mp3File), data, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	label := os.Args[1]
	out_dir := os.Args[2]
	SynthesizeSpeechInput(label, out_dir)
}
