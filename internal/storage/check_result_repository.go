package storage

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"health-monitor/internal/domain"
	"health-monitor/internal/storage/models"
)

// CheckResultRepository implements domain.CheckResultRepository
type CheckResultRepository struct {
	db *gorm.DB
}

// NewCheckResultRepository creates a new check result repository
func NewCheckResultRepository(db *gorm.DB) *CheckResultRepository {
	return &CheckResultRepository{db: db}
}

// Save saves a check result
func (r *CheckResultRepository) Save(ctx context.Context, result *domain.CheckResult) error {
	model := &models.CheckResult{}
	if err := model.FromDomain(result); err != nil {
		return fmt.Errorf("failed to convert result to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to save check result: %w", err)
	}

	// Update domain model with generated ID
	result.ID = model.ID

	return nil
}

// GetLatest retrieves the latest check result for a target
func (r *CheckResultRepository) GetLatest(ctx context.Context, targetID string) (*domain.CheckResult, error) {
	var model models.CheckResult
	if err := r.db.WithContext(ctx).
		Where("target_id = ?", targetID).
		Order("checked_at DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no check results found for target: %s", targetID)
		}
		return nil, fmt.Errorf("failed to get latest check result: %w", err)
	}

	result, err := model.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to result: %w", err)
	}

	return result, nil
}

// GetHistory retrieves check history for a target
func (r *CheckResultRepository) GetHistory(ctx context.Context, targetID string, limit, offset int) ([]*domain.CheckResult, error) {
	var modelResults []models.CheckResult
	query := r.db.WithContext(ctx).
		Where("target_id = ?", targetID).
		Order("checked_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&modelResults).Error; err != nil {
		return nil, fmt.Errorf("failed to get check history: %w", err)
	}

	results := make([]*domain.CheckResult, 0, len(modelResults))
	for _, model := range modelResults {
		result, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// GetHistoryInRange retrieves check history within a time range
func (r *CheckResultRepository) GetHistoryInRange(ctx context.Context, targetID string, from, to time.Time) ([]*domain.CheckResult, error) {
	var modelResults []models.CheckResult
	if err := r.db.WithContext(ctx).
		Where("target_id = ? AND checked_at BETWEEN ? AND ?", targetID, from, to).
		Order("checked_at DESC").
		Find(&modelResults).Error; err != nil {
		return nil, fmt.Errorf("failed to get check history in range: %w", err)
	}

	results := make([]*domain.CheckResult, 0, len(modelResults))
	for _, model := range modelResults {
		result, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// GetStats retrieves aggregated statistics for a target
func (r *CheckResultRepository) GetStats(ctx context.Context, targetID string, period time.Duration) (*domain.CheckStats, error) {
	since := time.Now().Add(-period)

	var stats struct {
		TotalChecks      int64
		SuccessfulChecks int64
		FailedChecks     int64
		AvgResponseTime  float64
		MinResponseTime  int64
		MaxResponseTime  int64
	}

	// Get aggregated statistics
	if err := r.db.WithContext(ctx).
		Model(&models.CheckResult{}).
		Select(`
			COUNT(*) as total_checks,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as successful_checks,
			SUM(CASE WHEN status != ? THEN 1 ELSE 0 END) as failed_checks,
			AVG(response_time_ms) as avg_response_time,
			MIN(response_time_ms) as min_response_time,
			MAX(response_time_ms) as max_response_time
		`, domain.CheckStatusSuccess, domain.CheckStatusSuccess).
		Where("target_id = ? AND checked_at >= ?", targetID, since).
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get check stats: %w", err)
	}

	// Get last check time
	var lastCheck models.CheckResult
	if err := r.db.WithContext(ctx).
		Where("target_id = ?", targetID).
		Order("checked_at DESC").
		First(&lastCheck).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get last check: %w", err)
	}

	// Get last success and failure times
	var lastSuccess models.CheckResult
	r.db.WithContext(ctx).
		Where("target_id = ? AND status = ?", targetID, domain.CheckStatusSuccess).
		Order("checked_at DESC").
		First(&lastSuccess)

	var lastFailure models.CheckResult
	r.db.WithContext(ctx).
		Where("target_id = ? AND status != ?", targetID, domain.CheckStatusSuccess).
		Order("checked_at DESC").
		First(&lastFailure)

	// Count consecutive failures
	var consecutiveFailures int
	rows, err := r.db.WithContext(ctx).
		Model(&models.CheckResult{}).
		Select("status").
		Where("target_id = ?", targetID).
		Order("checked_at DESC").
		Limit(100).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query consecutive failures: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			break
		}
		if status != string(domain.CheckStatusSuccess) {
			consecutiveFailures++
		} else {
			break
		}
	}

	// Calculate uptime percentage
	var uptimePercentage float64
	if stats.TotalChecks > 0 {
		uptimePercentage = (float64(stats.SuccessfulChecks) / float64(stats.TotalChecks)) * 100
	}

	// Determine current status
	currentStatus := domain.CheckStatusUnknown
	if lastCheck.ID > 0 {
		currentStatus = domain.CheckStatus(lastCheck.Status)
	}

	checkStats := &domain.CheckStats{
		TargetID:            targetID,
		TotalChecks:         stats.TotalChecks,
		SuccessfulChecks:    stats.SuccessfulChecks,
		FailedChecks:        stats.FailedChecks,
		UptimePercentage:    uptimePercentage,
		AvgResponseTimeMs:   stats.AvgResponseTime,
		MinResponseTimeMs:   stats.MinResponseTime,
		MaxResponseTimeMs:   stats.MaxResponseTime,
		LastCheckAt:         lastCheck.CheckedAt,
		CurrentStatus:       currentStatus,
		ConsecutiveFailures: consecutiveFailures,
		Period:              period,
	}

	if lastSuccess.ID > 0 {
		checkStats.LastSuccessAt = lastSuccess.CheckedAt
	}
	if lastFailure.ID > 0 {
		checkStats.LastFailureAt = lastFailure.CheckedAt
	}

	return checkStats, nil
}

// DeleteOlderThan deletes check results older than the specified time
func (r *CheckResultRepository) DeleteOlderThan(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("checked_at < ?", before).
		Delete(&models.CheckResult{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old check results: %w", result.Error)
	}

	return result.RowsAffected, nil
}
