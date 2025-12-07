package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type HomomorphicKeysRepository interface {
    Create(key *models.HomomorphicKey) error
    GetByID(id uint) (*models.HomomorphicKey, error)
    GetByElection(electionID uuid.UUID) (*models.HomomorphicKey, error)
    Update(key *models.HomomorphicKey) error
    Delete(id uint) error
}

type homomorphicKeysRepository struct {
    db *gorm.DB
}

func NewHomomorphicKeysRepository(db *gorm.DB) HomomorphicKeysRepository {
    return &homomorphicKeysRepository{db: db}
}

func (r *homomorphicKeysRepository) Create(key *models.HomomorphicKey) error {
    return r.db.Create(key).Error
}

func (r *homomorphicKeysRepository) GetByID(id uint) (*models.HomomorphicKey, error) {
    var key models.HomomorphicKey
    err := r.db.First(&key, id).Error
    return &key, err
}

func (r *homomorphicKeysRepository) GetByElection(electionID uuid.UUID) (*models.HomomorphicKey, error) {
    var key models.HomomorphicKey
    err := r.db.Where("elections = ?", electionID).First(&key).Error
    return &key, err
}

func (r *homomorphicKeysRepository) Update(key *models.HomomorphicKey) error {
    return r.db.Save(key).Error
}

func (r *homomorphicKeysRepository) Delete(id uint) error {
    return r.db.Delete(&models.HomomorphicKey{}, id).Error
}