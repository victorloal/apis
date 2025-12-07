package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type VotersService interface {
    CreateVoter(voter *models.Voter) error
    GetVoter(id uint) (*models.Voter, error)
    GetVotersByElection(electionID uuid.UUID) ([]models.Voter, error)
    GetVoterByToken(token string) (*models.Voter, error)
    UpdateVoter(voter *models.Voter) error
    DeleteVoter(id uint) error
    UpdateVoteStatus(id uint, status bool) error
}

type votersService struct {
    repo repositories.VotersRepository
}

func NewVotersService(repo repositories.VotersRepository) VotersService {
    return &votersService{repo: repo}
}

func (s *votersService) CreateVoter(voter *models.Voter) error {
    voter.CreatedAt = time.Now()
    voter.UpdatedAt = time.Now()
    return s.repo.Create(voter)
}

func (s *votersService) GetVoter(id uint) (*models.Voter, error) {
    return s.repo.GetByID(id)
}

func (s *votersService) GetVotersByElection(electionID uuid.UUID) ([]models.Voter, error) {
    return s.repo.GetByElection(electionID)
}

func (s *votersService) GetVoterByToken(token string) (*models.Voter, error) {
    return s.repo.GetByToken(token)
}

func (s *votersService) UpdateVoter(voter *models.Voter) error {
    voter.UpdatedAt = time.Now()
    return s.repo.Update(voter)
}

func (s *votersService) DeleteVoter(id uint) error {
    return s.repo.Delete(id)
}

func (s *votersService) UpdateVoteStatus(id uint, status bool) error {
    voter, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    
    voter.VoteStatus = status
    voter.UpdatedAt = time.Now()
    return s.repo.Update(voter)
}