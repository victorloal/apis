package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type ElectionService interface {
    CreateElection(election *models.Election) error
    GetElection(id uuid.UUID) (*models.Election, error)
    GetAllElections() ([]models.Election, error)
    UpdateElection(election *models.Election) error
    DeleteElection(id uuid.UUID) error
}

type electionService struct {
    repo repositories.ElectionRepository
}

func NewElectionService(repo repositories.ElectionRepository) ElectionService {
    return &electionService{repo: repo}
}

func (s *electionService) CreateElection(election *models.Election) error {
    election.CreatedAt = time.Now()
    election.UpdatedAt = time.Now()
    return s.repo.Create(election)
}

func (s *electionService) GetElection(id uuid.UUID) (*models.Election, error) {
    return s.repo.GetByID(id)
}

func (s *electionService) GetAllElections() ([]models.Election, error) {
    return s.repo.GetAll()
}

func (s *electionService) UpdateElection(election *models.Election) error {
    election.UpdatedAt = time.Now()
    return s.repo.Update(election)
}

func (s *electionService) DeleteElection(id uuid.UUID) error {
    return s.repo.Delete(id)
}