// pkg/cms/content.go
package cms

import (
    "context"
    "errors"
    "time"
)

type ContentService struct {
    db      *gorm.DB
    storage StorageService
    cache   CacheService
}

func (s *ContentService) Create(ctx context.Context, content *Content) error {
    if content.Slug == "" {
        content.Slug = generateSlug(content.Title)
    }
    
    return s.db.WithContext(ctx).Create(content).Error
}

func (s *ContentService) Get(ctx context.Context, id uint) (*Content, error) {
    var content Content
    
    // Try cache first
    if cached, err := s.cache.Get(cacheKey(id)); err == nil {
        return cached.(*Content), nil
    }
    
    if err := s.db.WithContext(ctx).First(&content, id).Error; err != nil {
        return nil, err
    }
    
    s.cache.Set(cacheKey(id), &content, time.Hour)
    return &content, nil
}

func (s *ContentService) List(ctx context.Context, filter ContentFilter) ([]Content, error) {
    var contents []Content
    query := s.db.WithContext(ctx)
    
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    if filter.AuthorID != 0 {
        query = query.Where("author_id = ?", filter.AuthorID)
    }
    
    return contents, query.Find(&contents).Error
}

func (s *ContentService) Update(ctx context.Context, content *Content) error {
    if err := s.db.WithContext(ctx).Save(content).Error; err != nil {
        return err
    }
    
    s.cache.Delete(cacheKey(content.ID))
    return nil
}

func (s *ContentService) Delete(ctx context.Context, id uint) error {
    if err := s.db.WithContext(ctx).Delete(&Content{}, id).Error; err != nil {
        return err
    }
    
    s.cache.Delete(cacheKey(id))
    return nil
}

func (s *ContentService) Publish(ctx context.Context, id uint) error {
    return s.db.WithContext(ctx).Model(&Content{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "status": "published",
            "published_at": time.Now(),
        }).Error
}

type ContentFilter struct {
    Status   string
    AuthorID uint
    Tags     []string
}

func cacheKey(id uint) string {
    return fmt.Sprintf("content:%d", id)
}
