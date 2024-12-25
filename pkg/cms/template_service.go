// pkg/cms/template_service.go
package cms

import (
    "context"
    "encoding/json"
    "errors"
    "time"
    "gorm.io/gorm"
)

type TemplateService struct {
    db    *gorm.DB
    cache CacheService
}

type Template struct {
    gorm.Model
    Name      string          `json:"name" gorm:"uniqueIndex"`
    Content   string          `json:"content"`
    Type      string          `json:"type"` // page, post, section
    Fields    json.RawMessage `json:"fields"`
    Version   int            `json:"version"`
    IsDefault bool           `json:"is_default"`
}

type TemplateField struct {
    Name     string `json:"name"`
    Type     string `json:"type"` // text, rich_text, image, etc.
    Required bool   `json:"required"`
}

type TemplateVersion struct {
    gorm.Model
    TemplateID uint   `json:"template_id"`
    Content    string `json:"content"`
    Version    int    `json:"version"`
}

func NewTemplateService(db *gorm.DB, cache CacheService) *TemplateService {
    return &TemplateService{
        db:    db,
        cache: cache,
    }
}

func (s *TemplateService) Create(ctx context.Context, template *Template) error {
    template.Version = 1
    err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(template).Error; err != nil {
            return err
        }

        version := &TemplateVersion{
            TemplateID: template.ID,
            Content:    template.Content,
            Version:    template.Version,
        }
        return tx.Create(version).Error
    })

    if err != nil {
        return err
    }

    s.cache.Delete("templates:list")
    return nil
}

func (s *TemplateService) Get(ctx context.Context, id uint) (*Template, error) {
    cacheKey := fmt.Sprintf("template:%d", id)
    
    if cached, err := s.cache.Get(cacheKey); err == nil {
        return cached.(*Template), nil
    }

    var template Template
    if err := s.db.WithContext(ctx).First(&template, id).Error; err != nil {
        return nil, err
    }

    s.cache.Set(cacheKey, &template, time.Hour)
    return &template, nil
}

func (s *TemplateService) Update(ctx context.Context, template *Template) error {
    template.Version++
    err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Save(template).Error; err != nil {
            return err
        }

        version := &TemplateVersion{
            TemplateID: template.ID,
            Content:    template.Content,
            Version:    template.Version,
        }
        return tx.Create(version).Error
    })

    if err != nil {
        return err
    }

    s.cache.Delete(fmt.Sprintf("template:%d", template.ID))
    s.cache.Delete("templates:list")
    return nil
}

func (s *TemplateService) Delete(ctx context.Context, id uint) error {
    err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Delete(&Template{}, id).Error; err != nil {
            return err
        }
        return tx.Where("template_id = ?", id).Delete(&TemplateVersion{}).Error
    })

    if err != nil {
        return err
    }

    s.cache.Delete(fmt.Sprintf("template:%d", id))
    s.cache.Delete("templates:list")
    return nil
}

func (s *TemplateService) GetVersion(ctx context.Context, templateID uint, version int) (*TemplateVersion, error) {
    var templateVersion TemplateVersion
    err := s.db.WithContext(ctx).
        Where("template_id = ? AND version = ?", templateID, version).
        First(&templateVersion).Error
    return &templateVersion, err
}

func (s *TemplateService) ListVersions(ctx context.Context, templateID uint) ([]TemplateVersion, error) {
    var versions []TemplateVersion
    err := s.db.WithContext(ctx).
        Where("template_id = ?", templateID).
        Order("version DESC").
        Find(&versions).Error
    return versions, err
}
