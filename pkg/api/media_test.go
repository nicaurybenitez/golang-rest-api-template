// pkg/cms/media_test.go
package cms

import (
    "bytes"
    "context"
    "io"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockS3Storage struct {
    mock.Mock
}

func (m *MockS3Storage) Upload(ctx context.Context, filePath string, content io.Reader, contentType string) (string, error) {
    args := m.Called(ctx, filePath, content, contentType)
    return args.String(0), args.Error(1)
}

func (m *MockS3Storage) Delete(ctx context.Context, key string) error {
    args := m.Called(ctx, key)
    return args.Error(0)
}

func (m *MockS3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
    args := m.Called(ctx, key)
    return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestMediaUpload(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
    assert.NoError(t, err)

    mockStorage := new(MockS3Storage)
    service := NewMediaService(gormDB, mockStorage)

    media := &Media{
        Name:     "test.jpg",
        MimeType: "image/jpeg",
        Size:     1024,
    }

    fileContent := []byte("test content")
    reader := bytes.NewReader(fileContent)

    expectedURL := "https://bucket.s3.region.amazonaws.com/test.jpg"
    mockStorage.On("Upload", mock.Anything, media.Name, reader, media.MimeType).Return(expectedURL, nil)

    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO media").WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    err = service.Upload(context.Background(), media, reader)
    assert.NoError(t, err)
    assert.Equal(t, expectedURL, media.URL)
}

func TestMediaDelete(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
    assert.NoError(t, err)

    mockStorage := new(MockS3Storage)
    service := NewMediaService(gormDB, mockStorage)

    media := &Media{
        Path: "test.jpg",
    }
    media.ID = 1

    rows := sqlmock.NewRows([]string{"id", "path"}).AddRow(1, "test.jpg")
    mock.ExpectQuery("^SELECT (.+) FROM media").WillReturnRows(rows)
    
    mockStorage.On("Delete", mock.Anything, media.Path).Return(nil)
    
    mock.ExpectBegin()
    mock.ExpectExec("DELETE FROM media").WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    err = service.Delete(context.Background(), media.ID)
    assert.NoError(t, err)
}

func TestMediaGet(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
    assert.NoError(t, err)

    mockStorage := new(MockS3Storage)
    service := NewMediaService(gormDB, mockStorage)

    rows := sqlmock.NewRows([]string{"id", "name", "type", "url"}).
        AddRow(1, "test.jpg", "image", "https://example.com/test.jpg")

    mock.ExpectQuery("^SELECT (.+) FROM media").WillReturnRows(rows)

    media, err := service.Get(context.Background(), 1)
    assert.NoError(t, err)
    assert.Equal(t, "test.jpg", media.Name)
}

func TestMediaList(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
    assert.NoError(t, err)

    mockStorage := new(MockS3Storage)
    service := NewMediaService(gormDB, mockStorage)

    rows := sqlmock.NewRows([]string{"id", "name", "type"}).
        AddRow(1, "test1.jpg", "image").
        AddRow(2, "test2.jpg", "image")

    mock.ExpectQuery("^SELECT (.+) FROM media").WillReturnRows(rows)

    filter := MediaFilter{Type: "image"}
    medias, err := service.List(context.Background(), filter)
    
    assert.NoError(t, err)
    assert.Equal(t, 2, len(medias))
}
