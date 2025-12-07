package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ElectionRepository interface {
    Create(election *models.Election) error
    GetByID(id uuid.UUID) (*models.Election, error)
    GetAll() ([]models.Election, error)
    Update(election *models.Election) error
    Delete(id uuid.UUID) error
}

type electionRepository struct {
    db *gorm.DB
}

func NewElectionRepository(db *gorm.DB) ElectionRepository {
    return &electionRepository{db: db}
}

func (r *electionRepository) Create(election *models.Election) error {
    return r.db.Create(election).Error
}

func (r *electionRepository) GetByID(id uuid.UUID) (*models.Election, error) {
    var election models.Election
    err := r.db.First(&election, "id = ?", id).Error
    return &election, err
}

func (r *electionRepository) GetAll() ([]models.Election, error) {
    var elections []models.Election
    err := r.db.Find(&elections).Error
    return elections, err
}

func (r *electionRepository) Update(election *models.Election) error {
    return r.db.Save(election).Error
}

func (r *electionRepository) Delete(id uuid.UUID) error {
    return r.db.Delete(&models.Election{}, "id = ?", id).Error
}