package services

import (
    "api_db/internal/models"
    "api_db/internal/repositories"
    "time"
)

type StatusService interface {
    CreateStatus(status *models.Status) error
    GetStatus(id int) (*models.Status, error)
    GetStatusByName(name string) (*models.Status, error)
    GetAllStatus() ([]models.Status, error)
    UpdateStatus(status *models.Status) error
    DeleteStatus(id int) error
}

type statusService struct {
    repo repositories.StatusRepository
}

func NewStatusService(repo repositories.StatusRepository) StatusService {
    return &statusService{repo: repo}
}

func (s *statusService) CreateStatus(status *models.Status) error {
    status.CreatedAt = time.Now()
    status.UpdatedAt = time.Now()
    return s.repo.Create(status)
}

func (s *statusService) GetStatus(id int) (*models.Status, error) {
    return s.repo.GetByID(id)
}

func (s *statusService) GetStatusByName(name string) (*models.Status, error) {
    return s.repo.GetByName(name)
}

func (s *statusService) GetAllStatus() ([]models.Status, error) {
    return s.repo.GetAll()
}

func (s *statusService) UpdateStatus(status *models.Status) error {
    status.UpdatedAt = time.Now()
    return s.repo.Update(status)
}

func (s *statusService) DeleteStatus(id int) error {
    return s.repo.Delete(id)
}