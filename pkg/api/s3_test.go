// pkg/storage/s3_test.go
package storage

import (
    "bytes"
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestS3Upload(t *testing.T) {
    // Mock AWS credentials for testing
    t.Setenv("AWS_ACCESS_KEY_ID", "test")
    t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
    t.Setenv("AWS_REGION", "us-east-1")

    storage, err := NewS3Storage("test-bucket")
    assert.NoError(t, err)

    content := bytes.NewReader([]byte("test content"))
    url, err := storage.Upload(context.Background(), "test.txt", content, "text/plain")
    assert.NoError(t, err)
    assert.Contains(t, url, "test-bucket")
}

func TestS3Delete(t *testing.T) {
    t.Setenv("AWS_ACCESS_KEY_ID", "test")
    t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
    t.Setenv("AWS_REGION", "us-east-1")

    storage, err := NewS3Storage("test-bucket")
    assert.NoError(t, err)

    err = storage.Delete(context.Background(), "test.txt")
    assert.NoError(t, err)
}
