package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	objectName   string
	objectBucket string
	objectType   string
	objectUri    string
	objectRegion string
	objectKeyId  string
	objectKey    string
)

const (
	helpText string = "" +
		"s3_setct\n" +
		"Set content type of an object in s3.\n" +
		"\n" +
		"Ex:\n" +
		"OBJECTBUCKET='bucketname'\n" +
		"OBJECTTYPE='application/epub+zip'\n" +
		"OBJECTURI='https://s3.us-west-002.backblazeb2.com'\n" +
		"OBJECTREGION='us-west-002'\n" +
		"OBJECTKEYID='<s3-keyId>'\n" +
		"OBJECTKEY='<s3-key>'\n" +
		" \n" +
		"echo 'some_ebook.epub' | s3_setct \n"
)

func readStdIn() ([]byte, error) {
	// check if there is somethinig to read on STDIN
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var stdin []byte
		scanner := bufio.NewScanner(os.Stdin)
		// max input size 10MB
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024*10)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return stdin, nil
	}
	return nil, errors.New("Expected stdin")
}

func loadEnv() error {
	objectBucket, _ = os.LookupEnv("OBJECTBUCKET")
	objectType, _ = os.LookupEnv("OBJECTTYPE")
	objectUri, _ = os.LookupEnv("OBJECTURI")
	objectRegion, _ = os.LookupEnv("OBJECTREGION")
	objectKeyId, _ = os.LookupEnv("OBJECTKEYID")
	objectKey, _ = os.LookupEnv("OBJECTKEY")

	for _, s := range []string{objectBucket, objectType, objectUri, objectRegion, objectKeyId, objectKey} {
		if s == "" {
			return fmt.Errorf("One or more environment variables not set")
		}
	}

	return nil
}

func main() {
	err := loadEnv()
	if err != nil {
		fmt.Print(helpText)
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	stdin, err := readStdIn()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Print(helpText)
		os.Exit(1)
	}

	objectName = strings.Trim(string(stdin), "\n ")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(objectKeyId, objectKey, ""),
		Endpoint:         aws.String(objectUri),
		Region:           aws.String(objectRegion),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	svc := s3.New(newSession)

	copyInput := &s3.CopyObjectInput{
		CopySource:        aws.String(objectBucket + "/" + objectName),
		Bucket:            aws.String(objectBucket), // dest bucket
		Key:               aws.String(objectName),   // dest object
		MetadataDirective: aws.String("REPLACE"),
		ContentType:       aws.String(objectType),
	}
	result, err := svc.CopyObject(copyInput)

	if err != nil {
		fmt.Printf("Failed to initiate update content type %s/%s -- %s, %s\n", objectBucket, objectName, objectType, err)
		os.Exit(1)
	}

	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(objectBucket), Key: aws.String(objectName)})
	if err != nil {
		fmt.Printf("Failed to complete update content type %s/%s -- %s, %s\n", objectBucket, objectName, objectType, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated content type %s -- %s\n%v\n", objectType, objectName, result)
	os.Exit(0)
}
