package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type TallyResultsService interface {
    CreateTallyResult(result *models.TallyResult) error
    GetTallyResult(id uint) (*models.TallyResult, error)
    GetTallyResultByElection(electionID uuid.UUID) (*models.TallyResult, error)
    UpdateTallyResult(result *models.TallyResult) error
    DeleteTallyResult(id uint) error
    ComputeTallyResult(electionID uuid.UUID, computedBy string) error
    GetTallyResultsWithElection() ([]map[string]interface{}, error)
}

type tallyResultsService struct {
    repo repositories.TallyResultsRepository
}

func NewTallyResultsService(repo repositories.TallyResultsRepository) TallyResultsService {
    return &tallyResultsService{repo: repo}
}

func (s *tallyResultsService) CreateTallyResult(result *models.TallyResult) error {
    result.CreatedAt = time.Now()
    result.UpdatedAt = time.Now()
    result.ComputedAt = time.Now()
    return s.repo.Create(result)
}

func (s *tallyResultsService) GetTallyResult(id uint) (*models.TallyResult, error) {
    return s.repo.GetByID(id)
}

func (s *tallyResultsService) GetTallyResultByElection(electionID uuid.UUID) (*models.TallyResult, error) {
    return s.repo.GetByElection(electionID)
}

func (s *tallyResultsService) UpdateTallyResult(result *models.TallyResult) error {
    result.UpdatedAt = time.Now()
    return s.repo.Update(result)
}

func (s *tallyResultsService) DeleteTallyResult(id uint) error {
    return s.repo.Delete(id)
}

func (s *tallyResultsService) ComputeTallyResult(electionID uuid.UUID, computedBy string) error {
    // Aquí iría la lógica para computar los resultados
    // Por ahora, creamos un resultado vacío
    result := &models.TallyResult{
        Election:   electionID,
        Results:    models.JSONB{},
        TotalVotes: 0,
        ComputedBy: computedBy,
        Proof:      models.JSONB{},
    }
    
    result.CreatedAt = time.Now()
    result.UpdatedAt = time.Now()
    result.ComputedAt = time.Now()
    
    return s.repo.Create(result)
}

func (s *tallyResultsService) GetTallyResultsWithElection() ([]map[string]interface{}, error) {
    return s.repo.GetWithElectionDetails()
}