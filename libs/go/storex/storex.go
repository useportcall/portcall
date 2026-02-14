package storex

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type IStore interface {
	GetFromSignatureBucket(id string, c context.Context) ([]byte, error)
	PutInSignatureBucket(id string, data []byte, c context.Context) error
	GetFromIconLogoBucket(id string, c context.Context) ([]byte, error)
	PutInIconLogoBucket(id string, data []byte, c context.Context) error
	DeleteFromIconLogoBucket(id string, c context.Context) error
}

type store struct {
	client *s3.Client
}

func New() (IStore, error) {
	return NewFromEnv()
}

func NewFromEnv() (IStore, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	region := os.Getenv("S3_REGION")
	accessKey := os.Getenv("S3_ACCESS_KEY_ID")
	secretKey := os.Getenv("S3_SECRET_ACCESS_KEY")

	// Log S3 configuration with credential presence only.
	log.Printf("S3 Config - Endpoint: %s, Region: %s, AccessKey: %s, SecretKey: %s",
		endpoint, region, func() string {
			if accessKey == "" {
				return "(empty)"
			}
			return "(set)"
		}(), func() string {
			if secretKey == "" {
				return "(empty)"
			}
			return "(set)"
		}())

	if accessKey == "" || secretKey == "" {
		log.Printf("WARNING: S3 credentials are not configured. File uploads will fail.")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	var s3Options []func(*s3.Options)
	if endpoint != "" {
		s3Options = append(s3Options, func(o *s3.Options) {
			o.BaseEndpoint = &endpoint
			o.UsePathStyle = true
		})
	} else {
		s3Options = append(s3Options, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	s3Client := s3.NewFromConfig(cfg, s3Options...)

	return &store{client: s3Client}, nil
}

func (s *store) GetFromSignatureBucket(id string, c context.Context) ([]byte, error) {
	bucket := "quote-signatures"
	key := fmt.Sprintf("%s.png", id)

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key}

	result, err := s.client.GetObject(c, input)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return bytes, nil
}

func (s *store) PutInSignatureBucket(id string, data []byte, c context.Context) error {
	bucket := "quote-signatures"
	key := fmt.Sprintf("%s.png", id)
	contentType := "image/png"

	input := &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType}

	_, err := s.client.PutObject(c, input)
	if err != nil {
		return err
	}

	return nil
}

//

func (s *store) GetFromIconLogoBucket(id string, c context.Context) ([]byte, error) {
	bucket := "icon-logos"
	key := fmt.Sprintf("%s.png", id)

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key}

	result, err := s.client.GetObject(c, input)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return bytes, nil
}

func (s *store) PutInIconLogoBucket(id string, data []byte, c context.Context) error {
	bucket := "icon-logos"
	key := fmt.Sprintf("%s.png", id)
	contentType := "image/png"

	input := &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	}

	_, err := s.client.PutObject(c, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) DeleteFromIconLogoBucket(id string, c context.Context) error {
	bucket := "icon-logos"
	key := fmt.Sprintf("%s.png", id)

	input := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	_, err := s.client.DeleteObject(c, input)
	if err != nil {
		return err
	}

	return nil
}
