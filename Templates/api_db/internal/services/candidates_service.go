package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type CandidatesService interface {
    CreateCandidate(candidate *models.Candidate) error
    GetCandidate(id int64) (*models.Candidate, error)
    GetCandidatesByElection(electionID uuid.UUID) ([]models.Candidate, error)
    UpdateCandidate(candidate *models.Candidate) error
    DeleteCandidate(id int64) error
    GetCandidatesByOrder(electionID uuid.UUID) ([]models.Candidate, error)
}

type candidatesService struct {
    repo repositories.CandidatesRepository
}

func NewCandidatesService(repo repositories.CandidatesRepository) CandidatesService {
    return &candidatesService{repo: repo}
}

func (s *candidatesService) CreateCandidate(candidate *models.Candidate) error {
    candidate.CreatedAt = time.Now()
    candidate.UpdateAt = time.Now()
    return s.repo.Create(candidate)
}

func (s *candidatesService) GetCandidate(id int64) (*models.Candidate, error) {
    return s.repo.GetByID(id)
}

func (s *candidatesService) GetCandidatesByElection(electionID uuid.UUID) ([]models.Candidate, error) {
    return s.repo.GetByElection(electionID)
}

func (s *candidatesService) UpdateCandidate(candidate *models.Candidate) error {
    candidate.UpdateAt = time.Now()
    return s.repo.Update(candidate)
}

func (s *candidatesService) DeleteCandidate(id int64) error {
    return s.repo.Delete(id)
}

func (s *candidatesService) GetCandidatesByOrder(electionID uuid.UUID) ([]models.Candidate, error) {
    return s.repo.GetByOrder(electionID)
}