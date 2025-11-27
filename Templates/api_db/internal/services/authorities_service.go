package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "github.com/google/uuid"
    "time"
)

type AuthoritiesService interface {
    CreateAuthority(authority *models.ElectionAuthority) error
    GetAuthority(id uint) (*models.ElectionAuthority, error)
    GetAuthoritiesByElection(electionID uuid.UUID) ([]models.ElectionAuthority, error)
    GetAuthorityByEmail(email string) (*models.ElectionAuthority, error)
    UpdateAuthority(authority *models.ElectionAuthority) error
    DeleteAuthority(id uint) error
}

type authoritiesService struct {
    repo repositories.AuthoritiesRepository
}

func NewAuthoritiesService(repo repositories.AuthoritiesRepository) AuthoritiesService {
    return &authoritiesService{repo: repo}
}

func (s *authoritiesService) CreateAuthority(authority *models.ElectionAuthority) error {
    authority.CreatedAt = time.Now()
    authority.UpdatedAt = time.Now()
    return s.repo.Create(authority)
}

func (s *authoritiesService) GetAuthority(id uint) (*models.ElectionAuthority, error) {
    return s.repo.GetByID(id)
}

func (s *authoritiesService) GetAuthoritiesByElection(electionID uuid.UUID) ([]models.ElectionAuthority, error) {
    return s.repo.GetByElection(electionID)
}

func (s *authoritiesService) GetAuthorityByEmail(email string) (*models.ElectionAuthority, error) {
    return s.repo.GetByEmail(email)
}

func (s *authoritiesService) UpdateAuthority(authority *models.ElectionAuthority) error {
    authority.UpdatedAt = time.Now()
    return s.repo.Update(authority)
}

func (s *authoritiesService) DeleteAuthority(id uint) error {
    return s.repo.Delete(id)
}