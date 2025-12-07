package repositories

import (
    "api_db/internal/models"
    "gorm.io/gorm"
)

type StatusRepository interface {
    Create(status *models.Status) error
    GetByID(id int) (*models.Status, error)
    GetByName(name string) (*models.Status, error)
    GetAll() ([]models.Status, error)
    Update(status *models.Status) error
    Delete(id int) error
}

type statusRepository struct {
    db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) StatusRepository {
    return &statusRepository{db: db}
}

func (r *statusRepository) Create(status *models.Status) error {
    return r.db.Create(status).Error
}

func (r *statusRepository) GetByID(id int) (*models.Status, error) {
    var status models.Status
    err := r.db.First(&status, id).Error
    return &status, err
}

func (r *statusRepository) GetByName(name string) (*models.Status, error) {
    var status models.Status
    err := r.db.Where("name = ?", name).First(&status).Error
    return &status, err
}

func (r *statusRepository) GetAll() ([]models.Status, error) {
    var statuses []models.Status
    err := r.db.Find(&statuses).Error
    return statuses, err
}

func (r *statusRepository) Update(status *models.Status) error {
    return r.db.Save(status).Error
}

func (r *statusRepository) Delete(id int) error {
    return r.db.Delete(&models.Status{}, id).Error
}