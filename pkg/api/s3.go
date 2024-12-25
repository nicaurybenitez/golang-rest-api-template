// pkg/storage/s3.go
package storage

import (
    "context"
    "fmt"
    "io"
    "path"
    "time"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type S3Storage struct {
    client *s3.Client
    bucket string
    uploader *manager.Uploader
}

func NewS3Storage(bucket string) (*S3Storage, error) {
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        return nil, fmt.Errorf("unable to load SDK config: %w", err)
    }

    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)

    return &S3Storage{
        client: client,
        bucket: bucket,
        uploader: uploader,
    }, nil
}

func (s *S3Storage) Upload(ctx context.Context, filePath string, content io.Reader, contentType string) (string, error) {
    key := fmt.Sprintf("%s/%d%s", 
        time.Now().Format("2006/01/02"),
        time.Now().UnixNano(),
        path.Ext(filePath),
    )

    result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        content,
        ContentType: aws.String(contentType),
    })
    if err != nil {
        return "", fmt.Errorf("failed to upload file: %w", err)
    }

    return result.Location, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
    _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return fmt.Errorf("failed to delete file: %w", err)
    }

    return nil
}

func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
    result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get file: %w", err)
    }

    return result.Body, nil
}
