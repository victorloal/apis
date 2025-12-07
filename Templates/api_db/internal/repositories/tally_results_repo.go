package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type TallyResultsRepository interface {
    Create(result *models.TallyResult) error
    GetByID(id uint) (*models.TallyResult, error)
    GetByElection(electionID uuid.UUID) (*models.TallyResult, error)
    Update(result *models.TallyResult) error
    Delete(id uint) error
    GetWithElectionDetails() ([]map[string]interface{}, error)
}

type tallyResultsRepository struct {
    db *gorm.DB
}

func NewTallyResultsRepository(db *gorm.DB) TallyResultsRepository {
    return &tallyResultsRepository{db: db}
}

func (r *tallyResultsRepository) Create(result *models.TallyResult) error {
    return r.db.Create(result).Error
}

func (r *tallyResultsRepository) GetByID(id uint) (*models.TallyResult, error) {
    var result models.TallyResult
    err := r.db.First(&result, id).Error
    return &result, err
}

func (r *tallyResultsRepository) GetByElection(electionID uuid.UUID) (*models.TallyResult, error) {
    var result models.TallyResult
    err := r.db.Where("election = ?", electionID).First(&result).Error
    return &result, err
}

func (r *tallyResultsRepository) Update(result *models.TallyResult) error {
    return r.db.Save(result).Error
}

func (r *tallyResultsRepository) Delete(id uint) error {
    return r.db.Delete(&models.TallyResult{}, id).Error
}

func (r *tallyResultsRepository) GetWithElectionDetails() ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    err := r.db.Table("tally_results").
        Select("tally_results.*, elections.name as election_name").
        Joins("LEFT JOIN elections ON tally_results.election = elections.id").
        Scan(&results).Error
    return results, err
}