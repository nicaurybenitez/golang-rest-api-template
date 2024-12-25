package storage

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
    client *s3.Client
    bucket string
}

func NewS3Client(bucket string) (*S3Client, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        return nil, err
    }

    return &S3Client{
        client: s3.NewFromConfig(cfg),
        bucket: bucket,
    }, nil
}
