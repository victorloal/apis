package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type CandidatesRepository interface {
    Create(candidate *models.Candidate) error
    GetByID(id int64) (*models.Candidate, error)
    GetByElection(electionID uuid.UUID) ([]models.Candidate, error)
    Update(candidate *models.Candidate) error
    Delete(id int64) error
    GetByOrder(electionID uuid.UUID) ([]models.Candidate, error)
}

type candidatesRepository struct {
    db *gorm.DB
}

func NewCandidatesRepository(db *gorm.DB) CandidatesRepository {
    return &candidatesRepository{db: db}
}

func (r *candidatesRepository) Create(candidate *models.Candidate) error {
    return r.db.Create(candidate).Error
}

func (r *candidatesRepository) GetByID(id int64) (*models.Candidate, error) {
    var candidate models.Candidate
    err := r.db.First(&candidate, id).Error
    return &candidate, err
}

func (r *candidatesRepository) GetByElection(electionID uuid.UUID) ([]models.Candidate, error) {
    var candidates []models.Candidate
    err := r.db.Where("elections = ?", electionID).Find(&candidates).Error
    return candidates, err
}

func (r *candidatesRepository) Update(candidate *models.Candidate) error {
    return r.db.Save(candidate).Error
}

func (r *candidatesRepository) Delete(id int64) error {
    return r.db.Delete(&models.Candidate{}, id).Error
}

func (r *candidatesRepository) GetByOrder(electionID uuid.UUID) ([]models.Candidate, error) {
    var candidates []models.Candidate
    err := r.db.Where("elections = ?", electionID).Order("candidate_order").Find(&candidates).Error
    return candidates, err
}