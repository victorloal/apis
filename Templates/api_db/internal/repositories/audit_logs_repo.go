package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "time"
)

type AuditLogsRepository interface {
    Create(log *models.AuditLog) error
    GetByID(id uint) (*models.AuditLog, error)
    GetByElection(electionID uuid.UUID) ([]models.AuditLog, error)
    GetByAction(action string) ([]models.AuditLog, error)
    GetByUser(userType, userID string) ([]models.AuditLog, error)
    GetByDateRange(start, end time.Time) ([]models.AuditLog, error)
    Delete(id uint) error
}

type auditLogsRepository struct {
    db *gorm.DB
}

func NewAuditLogsRepository(db *gorm.DB) AuditLogsRepository {
    return &auditLogsRepository{db: db}
}

func (r *auditLogsRepository) Create(log *models.AuditLog) error {
    return r.db.Create(log).Error
}

func (r *auditLogsRepository) GetByID(id uint) (*models.AuditLog, error) {
    var log models.AuditLog
    err := r.db.First(&log, id).Error
    return &log, err
}

func (r *auditLogsRepository) GetByElection(electionID uuid.UUID) ([]models.AuditLog, error) {
    var logs []models.AuditLog
    err := r.db.Where("election = ?", electionID).Order("created_at DESC").Find(&logs).Error
    return logs, err
}

func (r *auditLogsRepository) GetByAction(action string) ([]models.AuditLog, error) {
    var logs []models.AuditLog
    err := r.db.Where("action = ?", action).Order("created_at DESC").Find(&logs).Error
    return logs, err
}

func (r *auditLogsRepository) GetByUser(userType, userID string) ([]models.AuditLog, error) {
    var logs []models.AuditLog
    err := r.db.Where("user_type = ? AND user_id = ?", userType, userID).Order("created_at DESC").Find(&logs).Error
    return logs, err
}

func (r *auditLogsRepository) GetByDateRange(start, end time.Time) ([]models.AuditLog, error) {
    var logs []models.AuditLog
    err := r.db.Where("created_at BETWEEN ? AND ?", start, end).Order("created_at DESC").Find(&logs).Error
    return logs, err
}

func (r *auditLogsRepository) Delete(id uint) error {
    return r.db.Delete(&models.AuditLog{}, id).Error
}