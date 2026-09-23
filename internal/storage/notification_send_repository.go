package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"health-monitor/internal/domain"
	"health-monitor/internal/storage/models"
)

// NotificationSendRepository implements domain.NotificationSendRepository
type NotificationSendRepository struct {
	db *gorm.DB
}

// NewNotificationSendRepository creates a new notification send repository
func NewNotificationSendRepository(db *gorm.DB) *NotificationSendRepository {
	return &NotificationSendRepository{db: db}
}

// Claim inserts a send row, relying on the unique index to reject duplicates.
// DoNothing is used instead of a driver error because the GORM connection does
// not enable TranslateError, so gorm.ErrDuplicatedKey would never match.
func (r *NotificationSendRepository) Claim(ctx context.Context, send *domain.NotificationSend) error {
	model := &models.NotificationSend{}
	model.FromDomain(send)

	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(model)

	if result.Error != nil {
		return fmt.Errorf("failed to claim notification send: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrSendAlreadyClaimed
	}

	send.ID = model.ID
	send.CreatedAt = model.CreatedAt

	return nil
}

// Complete stores the terminal state of a previously claimed send
func (r *NotificationSendRepository) Complete(ctx context.Context, send *domain.NotificationSend) error {
	result := r.db.WithContext(ctx).
		Model(&models.NotificationSend{}).
		Where("id = ?", send.ID).
		Updates(map[string]interface{}{
			"status":          string(send.Status),
			"external_id":     send.ExternalID,
			"external_status": send.ExternalStatus,
			"cost":            send.Cost,
			"error":           send.Error,
			"completed_at":    send.CompletedAt,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to complete notification send: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification send not found: %d", send.ID)
	}

	return nil
}

// ListByIncident retrieves sends of one notifier for one incident, oldest first
func (r *NotificationSendRepository) ListByIncident(ctx context.Context, incidentID int64, notifierID string) ([]*domain.NotificationSend, error) {
	var found []models.NotificationSend

	if err := r.db.WithContext(ctx).
		Where("incident_id = ? AND notifier_id = ?", incidentID, notifierID).
		Order("sequence ASC").
		Find(&found).Error; err != nil {
		return nil, fmt.Errorf("failed to list notification sends: %w", err)
	}

	sends := make([]*domain.NotificationSend, 0, len(found))
	for i := range found {
		sends = append(sends, found[i].ToDomain())
	}

	return sends, nil
}
