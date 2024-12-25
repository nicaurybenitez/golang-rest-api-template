// pkg/cms/scheduler_service.go
package cms

import (
    "context"
    "time"
    "gorm.io/gorm"
)

type SchedulerService struct {
    db      *gorm.DB
    content *ContentService
}

type ScheduledPublication struct {
    gorm.Model
    ContentID    uint      `json:"content_id"`
    PublishAt    time.Time `json:"publish_at"`
    Status       string    `json:"status"` // pending, completed, failed
    ErrorMessage string    `json:"error_message,omitempty"`
}

func NewSchedulerService(db *gorm.DB, content *ContentService) *SchedulerService {
    return &SchedulerService{
        db:      db,
        content: content,
    }
}

func (s *SchedulerService) Schedule(ctx context.Context, contentID uint, publishAt time.Time) error {
    schedule := &ScheduledPublication{
        ContentID: contentID,
        PublishAt: publishAt,
        Status:    "pending",
    }
    return s.db.WithContext(ctx).Create(schedule).Error
}

func (s *SchedulerService) ProcessScheduled(ctx context.Context) error {
    var schedules []ScheduledPublication
    err := s.db.WithContext(ctx).
        Where("status = ? AND publish_at <= ?", "pending", time.Now()).
        Find(&schedules).Error
    if err != nil {
        return err
    }

    for _, schedule := range schedules {
        err := s.content.Publish(ctx, schedule.ContentID)
        if err != nil {
            s.updateStatus(ctx, schedule.ID, "failed", err.Error())
            continue
        }
        s.updateStatus(ctx, schedule.ID, "completed", "")
    }
    return nil
}

func (s *SchedulerService) updateStatus(ctx context.Context, id uint, status, errorMsg string) error {
    return s.db.WithContext(ctx).Model(&ScheduledPublication{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "status":        status,
            "error_message": errorMsg,
        }).Error
}

func (s *SchedulerService) Cancel(ctx context.Context, scheduleID uint) error {
    return s.db.WithContext(ctx).Delete(&ScheduledPublication{}, scheduleID).Error
}

func (s *SchedulerService) List(ctx context.Context, status string) ([]ScheduledPublication, error) {
    var schedules []ScheduledPublication
    query := s.db.WithContext(ctx)
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    return schedules, query.Find(&schedules).Error
}
