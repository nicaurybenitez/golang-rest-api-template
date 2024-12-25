// pkg/cms/media.go
package cms

import (
    "context"
    "io"
    "path/filepath"
    "gorm.io/gorm"
    "your-project/pkg/storage"
)

type MediaService struct {
    db      *gorm.DB
    storage *storage.S3Storage
}

func NewMediaService(db *gorm.DB, storage *storage.S3Storage) *MediaService {
    return &MediaService{
        db:      db,
        storage: storage,
    }
}

func (s *MediaService) Upload(ctx context.Context, media *Media, file io.Reader) error {
    url, err := s.storage.Upload(ctx, media.Name, file, media.MimeType)
    if err != nil {
        return err
    }

    media.URL = url
    media.Path = filepath.Base(url)

    return s.db.WithContext(ctx).Create(media).Error
}

func (s *MediaService) Get(ctx context.Context, id uint) (*Media, error) {
    var media Media
    if err := s.db.WithContext(ctx).First(&media, id).Error; err != nil {
        return nil, err
    }
    return &media, nil
}

func (s *MediaService) List(ctx context.Context, filter MediaFilter) ([]Media, error) {
    var medias []Media
    query := s.db.WithContext(ctx)

    if filter.Type != "" {
        query = query.Where("type = ?", filter.Type)
    }

    return medias, query.Find(&medias).Error
}

func (s *MediaService) Delete(ctx context.Context, id uint) error {
    var media Media
    if err := s.db.WithContext(ctx).First(&media, id).Error; err != nil {
        return err
    }

    if err := s.storage.Delete(ctx, media.Path); err != nil {
        return err
    }

    return s.db.WithContext(ctx).Delete(&media).Error
}

type Media struct {
    gorm.Model
    Name     string `json:"name"`
    Type     string `json:"type"`
    URL      string `json:"url"`
    Size     int64  `json:"size"`
    Path     string `json:"path"`
    MimeType string `json:"mime_type"`
}

type MediaFilter struct {
    Type string
}
