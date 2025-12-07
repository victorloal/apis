package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type BallotsRepository interface {
    Create(ballot *models.Ballot) error
    GetByID(id string, electionID uuid.UUID, voterID int) (*models.Ballot, error)
    GetByElection(electionID uuid.UUID) ([]models.Ballot, error)
    GetByVoter(voterID int) ([]models.Ballot, error)
    Update(ballot *models.Ballot) error
    Delete(id string, electionID uuid.UUID, voterID int) error
    GetWithVoterDetails(electionID uuid.UUID) ([]map[string]interface{}, error)
}

type ballotsRepository struct {
    db *gorm.DB
}

func NewBallotsRepository(db *gorm.DB) BallotsRepository {
    return &ballotsRepository{db: db}
}

func (r *ballotsRepository) Create(ballot *models.Ballot) error {
    return r.db.Create(ballot).Error
}

func (r *ballotsRepository) GetByID(id string, electionID uuid.UUID, voterID int) (*models.Ballot, error) {
    var ballot models.Ballot
    err := r.db.Where("id = ? AND elections = ? AND voter = ?", id, electionID, voterID).First(&ballot).Error
    return &ballot, err
}

func (r *ballotsRepository) GetByElection(electionID uuid.UUID) ([]models.Ballot, error) {
    var ballots []models.Ballot
    err := r.db.Where("elections = ?", electionID).Find(&ballots).Error
    return ballots, err
}

func (r *ballotsRepository) GetByVoter(voterID int) ([]models.Ballot, error) {
    var ballots []models.Ballot
    err := r.db.Where("voter = ?", voterID).Find(&ballots).Error
    return ballots, err
}

func (r *ballotsRepository) Update(ballot *models.Ballot) error {
    return r.db.Save(ballot).Error
}

func (r *ballotsRepository) Delete(id string, electionID uuid.UUID, voterID int) error {
    return r.db.Where("id = ? AND elections = ? AND voter = ?", id, electionID, voterID).Delete(&models.Ballot{}).Error
}

func (r *ballotsRepository) GetWithVoterDetails(electionID uuid.UUID) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    err := r.db.Table("ballots").
        Select("ballots.*, voters.token, voters.vote_status").
        Joins("LEFT JOIN voters ON ballots.voter = voters.id").
        Where("ballots.elections = ?", electionID).
        Scan(&results).Error
    return results, err
}