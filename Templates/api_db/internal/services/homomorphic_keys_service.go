package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type HomomorphicKeysService interface {
    CreateKey(key *models.HomomorphicKey) error
    GetKey(id uint) (*models.HomomorphicKey, error)
    GetKeyByElection(electionID uuid.UUID) (*models.HomomorphicKey, error)
    UpdateKey(key *models.HomomorphicKey) error
    DeleteKey(id uint) error
    UpdateKeyParams(electionID uuid.UUID, params models.JSONB) error
}

type homomorphicKeysService struct {
    repo repositories.HomomorphicKeysRepository
}

func NewHomomorphicKeysService(repo repositories.HomomorphicKeysRepository) HomomorphicKeysService {
    return &homomorphicKeysService{repo: repo}
}

func (s *homomorphicKeysService) CreateKey(key *models.HomomorphicKey) error {
    key.CreatedAt = time.Now()
    key.UpdatedAt = time.Now()
    return s.repo.Create(key)
}

func (s *homomorphicKeysService) GetKey(id uint) (*models.HomomorphicKey, error) {
    return s.repo.GetByID(id)
}

func (s *homomorphicKeysService) GetKeyByElection(electionID uuid.UUID) (*models.HomomorphicKey, error) {
    return s.repo.GetByElection(electionID)
}

func (s *homomorphicKeysService) UpdateKey(key *models.HomomorphicKey) error {
    key.UpdatedAt = time.Now()
    return s.repo.Update(key)
}

func (s *homomorphicKeysService) DeleteKey(id uint) error {
    return s.repo.Delete(id)
}

func (s *homomorphicKeysService) UpdateKeyParams(electionID uuid.UUID, params models.JSONB) error {
    key, err := s.repo.GetByElection(electionID)
    if err != nil {
        return err
    }
    
    key.Params = params
    key.UpdatedAt = time.Now()
    return s.repo.Update(key)
}