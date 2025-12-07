package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type AuditConfigRepository interface {
    Create(config *models.ElectionAuditConfig) error
    GetByID(id uint) (*models.ElectionAuditConfig, error)
    GetByElection(electionID uuid.UUID) (*models.ElectionAuditConfig, error)
    Update(config *models.ElectionAuditConfig) error
    Delete(id uint) error
}

type auditConfigRepository struct {
    db *gorm.DB
}

func NewAuditConfigRepository(db *gorm.DB) AuditConfigRepository {
    return &auditConfigRepository{db: db}
}

func (r *auditConfigRepository) Create(config *models.ElectionAuditConfig) error {
    return r.db.Create(config).Error
}

func (r *auditConfigRepository) GetByID(id uint) (*models.ElectionAuditConfig, error) {
    var config models.ElectionAuditConfig
    err := r.db.First(&config, id).Error
    return &config, err
}

func (r *auditConfigRepository) GetByElection(electionID uuid.UUID) (*models.ElectionAuditConfig, error) {
    var config models.ElectionAuditConfig
    err := r.db.Where("election = ?", electionID).First(&config).Error
    return &config, err
}

func (r *auditConfigRepository) Update(config *models.ElectionAuditConfig) error {
    return r.db.Save(config).Error
}

func (r *auditConfigRepository) Delete(id uint) error {
    return r.db.Delete(&models.ElectionAuditConfig{}, id).Error
}