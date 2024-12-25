// pkg/cms/version_service.go
package cms

import (
    "context"
    "time"
    "gorm.io/gorm"
)

type VersionService struct {
    db *gorm.DB
}

type Version struct {
    gorm.Model
    EntityID   uint        `json:"entity_id"`
    EntityType string      `json:"entity_type"` // content, template
    Data       JSON        `json:"data"`
    Version    int         `json:"version"`
    CreatedBy  uint        `json:"created_by"`
    Comment    string      `json:"comment"`
}

func NewVersionService(db *gorm.DB) *VersionService {
    return &VersionService{db: db}
}

func (s *VersionService) CreateVersion(ctx context.Context, entityType string, entityID uint, data interface{}, createdBy uint, comment string) error {
    lastVersion, _ := s.GetLastVersion(ctx, entityType, entityID)
    version := &Version{
        EntityID:   entityID,
        EntityType: entityType,
        Data:       data,
        Version:    lastVersion + 1,
        CreatedBy:  createdBy,
        Comment:    comment,
    }
    
    return s.db.WithContext(ctx).Create(version).Error
}

func (s *VersionService) GetVersion(ctx context.Context, entityType string, entityID uint, version int) (*Version, error) {
    var v Version
    err := s.db.WithContext(ctx).
        Where("entity_type = ? AND entity_id = ? AND version = ?", entityType, entityID, version).
        First(&v).Error
    return &v, err
}

func (s *VersionService) GetLastVersion(ctx context.Context, entityType string, entityID uint) (int, error) {
    var v Version
    err := s.db.WithContext(ctx).
        Where("entity_type = ? AND entity_id = ?", entityType, entityID).
        Order("version desc").
        First(&v).Error
    
    if err == gorm.ErrRecordNotFound {
        return 0, nil
    }
    return v.Version, err
}

func (s *VersionService) ListVersions(ctx context.Context, entityType string, entityID uint) ([]Version, error) {
    var versions []Version
    err := s.db.WithContext(ctx).
        Where("entity_type = ? AND entity_id = ?", entityType, entityID).
        Order("version desc").
        Find(&versions).Error
    return versions, err
}

func (s *VersionService) Restore(ctx context.Context, entityType string, entityID uint, version int) error {
    v, err := s.GetVersion(ctx, entityType, entityID, version)
    if err != nil {
        return err
    }

    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        switch entityType {
        case "content":
            var content Content
            if err := json.Unmarshal(v.Data, &content); err != nil {
                return err
            }
            return tx.Model(&Content{}).Where("id = ?", entityID).Updates(content).Error
        case "template":
            var template Template
            if err := json.Unmarshal(v.Data, &template); err != nil {
                return err
            }
            return tx.Model(&Template{}).Where("id = ?", entityID).Updates(template).Error
        }
        return nil
    })
}

func (s *VersionService) Compare(ctx context.Context, entityType string, entityID uint, v1, v2 int) (map[string]interface{}, error) {
    version1, err := s.GetVersion(ctx, entityType, entityID, v1)
    if err != nil {
        return nil, err
    }

    version2, err := s.GetVersion(ctx, entityType, entityID, v2)
    if err != nil {
        return nil, err
    }

    var data1, data2 map[string]interface{}
    json.Unmarshal(version1.Data, &data1)
    json.Unmarshal(version2.Data, &data2)

    diff := make(map[string]interface{})
    for k, v := range data1 {
        if data2[k] != v {
            diff[k] = map[string]interface{}{
                "old": v,
                "new": data2[k],
            }
        }
    }

    return diff, nil
}
