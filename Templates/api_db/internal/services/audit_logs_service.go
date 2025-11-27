package services

import (
	"api_db/internal/models"
	"api_db/internal/repositories"
	"net"
	"time"

	"github.com/google/uuid"
)

type AuditLogsService interface {
    CreateAuditLog(log *models.AuditLog) error
    GetAuditLog(id uint) (*models.AuditLog, error)
    GetAuditLogsByElection(electionID uuid.UUID) ([]models.AuditLog, error)
    GetAuditLogsByAction(action string) ([]models.AuditLog, error)
    GetAuditLogsByUser(userType, userID string) ([]models.AuditLog, error)
    GetAuditLogsByDateRange(start, end time.Time) ([]models.AuditLog, error)
    DeleteAuditLog(id uint) error
    LogVoteAction(electionID uuid.UUID, voterID string, ipAddress string, userAgent string) error
    LogAuthorityAction(electionID uuid.UUID, authorityID string, action string, details models.JSONB) error
}

type auditLogsService struct {
    repo repositories.AuditLogsRepository
}

func NewAuditLogsService(repo repositories.AuditLogsRepository) AuditLogsService {
    return &auditLogsService{repo: repo}
}

func (s *auditLogsService) CreateAuditLog(log *models.AuditLog) error {
    log.CreatedAt = time.Now()
    return s.repo.Create(log)
}

func (s *auditLogsService) GetAuditLog(id uint) (*models.AuditLog, error) {
    return s.repo.GetByID(id)
}

func (s *auditLogsService) GetAuditLogsByElection(electionID uuid.UUID) ([]models.AuditLog, error) {
    return s.repo.GetByElection(electionID)
}

func (s *auditLogsService) GetAuditLogsByAction(action string) ([]models.AuditLog, error) {
    return s.repo.GetByAction(action)
}

func (s *auditLogsService) GetAuditLogsByUser(userType, userID string) ([]models.AuditLog, error) {
    return s.repo.GetByUser(userType, userID)
}

func (s *auditLogsService) GetAuditLogsByDateRange(start, end time.Time) ([]models.AuditLog, error) {
    return s.repo.GetByDateRange(start, end)
}

func (s *auditLogsService) DeleteAuditLog(id uint) error {
    return s.repo.Delete(id)
}

func (s *auditLogsService) LogVoteAction(electionID uuid.UUID, voterID string, ipAddress string, userAgent string) error {
    log := &models.AuditLog{
        Election:  electionID,
        Action:    "vote_cast",
        UserType:  "voter",
        UserID:    voterID,
        IPAddress: net.ParseIP(ipAddress),
        UserAgent: userAgent,
        Details:   models.JSONB{"action": "vote_cast", "timestamp": time.Now()},
    }
    return s.CreateAuditLog(log)
}

func (s *auditLogsService) LogAuthorityAction(electionID uuid.UUID, authorityID string, action string, details models.JSONB) error {
    log := &models.AuditLog{
        Election:  electionID,
        Action:    action,
        UserType:  "authority",
        UserID:    authorityID,
        Details:   details,
    }
    return s.CreateAuditLog(log)
}