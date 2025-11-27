package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type VotersRepository interface {
    Create(voter *models.Voter) error
    GetByID(id uint) (*models.Voter, error)
    GetByElection(electionID uuid.UUID) ([]models.Voter, error)
    GetByToken(token string) (*models.Voter, error)
    Update(voter *models.Voter) error
    Delete(id uint) error
}

type votersRepository struct {
    db *gorm.DB
}

func NewVotersRepository(db *gorm.DB) VotersRepository {
    return &votersRepository{db: db}
}

func (r *votersRepository) Create(voter *models.Voter) error {
    return r.db.Create(voter).Error
}

func (r *votersRepository) GetByID(id uint) (*models.Voter, error) {
    var voter models.Voter
    err := r.db.First(&voter, id).Error
    return &voter, err
}

func (r *votersRepository) GetByElection(electionID uuid.UUID) ([]models.Voter, error) {
    var voters []models.Voter
    err := r.db.Where("elections = ?", electionID).Find(&voters).Error
    return voters, err
}

func (r *votersRepository) GetByToken(token string) (*models.Voter, error) {
    var voter models.Voter
    err := r.db.Where("token = ?", token).First(&voter).Error
    return &voter, err
}

func (r *votersRepository) Update(voter *models.Voter) error {
    return r.db.Save(voter).Error
}

func (r *votersRepository) Delete(id uint) error {
    return r.db.Delete(&models.Voter{}, id).Error
}