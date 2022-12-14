package main

import (
	"context"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/dustin/go-humanize"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
)

// PutObject - Upload object to s3 bucket
func PutObject(key, bucket, s3Class string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	session := s3.NewFromConfig(cfg)

	file, err := os.Open(key)
	if err != nil {
		return err
	}
	defer file.Close()

	i := &s3.PutObjectInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(key),
		Body:         file,
		StorageClass: types.StorageClass(s3Class),
	}

	_, err = session.PutObject(context.TODO(), i)
	if err == nil {
		log.Print("Cache saved successfully")
	}

	return err
}

// GetObject - Get object from s3 bucket
func GetObject(key, bucket string) error {

	downloadedFile, err := os.Create(key)
	if err != nil {
		log.Println(err)
	}
	defer downloadedFile.Close()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	session := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	downloader := manager.NewDownloader(session)
	numBytes, err := downloader.Download(context.TODO(), downloadedFile, input)
	if err == nil {
		log.Printf("Cache downloaded successfully, size %s", humanize.Bytes(uint64(numBytes)))
	}

	return err
}

// DeleteObject - Delete object from s3 bucket
func DeleteObject(key, bucket string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	session := s3.NewFromConfig(cfg)

	i := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err = session.DeleteObject(context.TODO(), i)
	if err == nil {
		log.Print("Cache purged successfully")
	}

	return err
}

// ObjectExists - Verify if object exists in s3
func ObjectExists(key, bucket string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	session := s3.NewFromConfig(cfg)

	i := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err = session.HeadObject(context.TODO(), i); err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
