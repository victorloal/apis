package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type BallotsService interface {
    CreateBallot(ballot *models.Ballot) error
    GetBallot(id string, electionID uuid.UUID, voterID int) (*models.Ballot, error)
    GetBallotsByElection(electionID uuid.UUID) ([]models.Ballot, error)
    GetBallotsByVoter(voterID int) ([]models.Ballot, error)
    UpdateBallot(ballot *models.Ballot) error
    DeleteBallot(id string, electionID uuid.UUID, voterID int) error
    GetBallotsWithVoterDetails(electionID uuid.UUID) ([]map[string]interface{}, error)
}

type ballotsService struct {
    repo repositories.BallotsRepository
}

func NewBallotsService(repo repositories.BallotsRepository) BallotsService {
    return &ballotsService{repo: repo}
}

func (s *ballotsService) CreateBallot(ballot *models.Ballot) error {
    ballot.CreatedAt = time.Now()
    ballot.UpdatedAt = time.Now()
    return s.repo.Create(ballot)
}

func (s *ballotsService) GetBallot(id string, electionID uuid.UUID, voterID int) (*models.Ballot, error) {
    return s.repo.GetByID(id, electionID, voterID)
}

func (s *ballotsService) GetBallotsByElection(electionID uuid.UUID) ([]models.Ballot, error) {
    return s.repo.GetByElection(electionID)
}

func (s *ballotsService) GetBallotsByVoter(voterID int) ([]models.Ballot, error) {
    return s.repo.GetByVoter(voterID)
}

func (s *ballotsService) UpdateBallot(ballot *models.Ballot) error {
    ballot.UpdatedAt = time.Now()
    return s.repo.Update(ballot)
}

func (s *ballotsService) DeleteBallot(id string, electionID uuid.UUID, voterID int) error {
    return s.repo.Delete(id, electionID, voterID)
}

func (s *ballotsService) GetBallotsWithVoterDetails(electionID uuid.UUID) ([]map[string]interface{}, error) {
    return s.repo.GetWithVoterDetails(electionID)
}