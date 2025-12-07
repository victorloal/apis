package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type AuditConfigService interface {
    CreateAuditConfig(config *models.ElectionAuditConfig) error
    GetAuditConfig(id uint) (*models.ElectionAuditConfig, error)
    GetAuditConfigByElection(electionID uuid.UUID) (*models.ElectionAuditConfig, error)
    UpdateAuditConfig(config *models.ElectionAuditConfig) error
    DeleteAuditConfig(id uint) error
    EnableBallotAudit(electionID uuid.UUID, enable bool) error
    EnableAccessLogs(electionID uuid.UUID, enable bool) error
}

type auditConfigService struct {
    repo repositories.AuditConfigRepository
}

func NewAuditConfigService(repo repositories.AuditConfigRepository) AuditConfigService {
    return &auditConfigService{repo: repo}
}

func (s *auditConfigService) CreateAuditConfig(config *models.ElectionAuditConfig) error {
    config.CreatedAt = time.Now()
    config.UpdatedAt = time.Now()
    return s.repo.Create(config)
}

func (s *auditConfigService) GetAuditConfig(id uint) (*models.ElectionAuditConfig, error) {
    return s.repo.GetByID(id)
}

func (s *auditConfigService) GetAuditConfigByElection(electionID uuid.UUID) (*models.ElectionAuditConfig, error) {
    return s.repo.GetByElection(electionID)
}

func (s *auditConfigService) UpdateAuditConfig(config *models.ElectionAuditConfig) error {
    config.UpdatedAt = time.Now()
    return s.repo.Update(config)
}

func (s *auditConfigService) DeleteAuditConfig(id uint) error {
    return s.repo.Delete(id)
}

func (s *auditConfigService) EnableBallotAudit(electionID uuid.UUID, enable bool) error {
    config, err := s.repo.GetByElection(electionID)
    if err != nil {
        return err
    }
    
    config.EnableBallotAudit = enable
    config.UpdatedAt = time.Now()
    return s.repo.Update(config)
}

func (s *auditConfigService) EnableAccessLogs(electionID uuid.UUID, enable bool) error {
    config, err := s.repo.GetByElection(electionID)
    if err != nil {
        return err
    }
    
    config.EnableAccessLogs = enable
    config.UpdatedAt = time.Now()
    return s.repo.Update(config)
}