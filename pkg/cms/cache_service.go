// pkg/cms/cache_service.go
package cms

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type CacheService struct {
    client *redis.Client
}

func NewCacheService(client *redis.Client) *CacheService {
    return &CacheService{
        client: client,
    }
}

func (s *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return s.client.Set(context.Background(), key, data, ttl).Err()
}

func (s *CacheService) Get(key string, result interface{}) error {
    data, err := s.client.Get(context.Background(), key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, result)
}

func (s *CacheService) Delete(key string) error {
    return s.client.Del(context.Background(), key).Err()
}

func (s *CacheService) DeletePattern(pattern string) error {
    iter := s.client.Scan(context.Background(), 0, pattern, 0).Iterator()
    for iter.Next(context.Background()) {
        if err := s.Delete(iter.Val()); err != nil {
            return err
        }
    }
    return iter.Err()
}

// Métodos específicos para CMS
func (s *CacheService) CacheContent(content *Content) error {
    key := fmt.Sprintf("content:%d", content.ID)
    return s.Set(key, content, time.Hour)
}

func (s *CacheService) CacheTemplate(template *Template) error {
    key := fmt.Sprintf("template:%d", template.ID)
    return s.Set(key, template, time.Hour)
}

func (s *CacheService) InvalidateContent(contentID uint) error {
    key := fmt.Sprintf("content:%d", contentID)
    return s.Delete(key)
}

func (s *CacheService) InvalidateTemplate(templateID uint) error {
    key := fmt.Sprintf("template:%d", templateID)
    if err := s.Delete(key); err != nil {
        return err
    }
    // Invalidar también el contenido que usa esta plantilla
    return s.DeletePattern("content:*")
}
