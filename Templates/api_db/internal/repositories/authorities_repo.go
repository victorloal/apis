package repositories

import (
    "api_db/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type AuthoritiesRepository interface {
    Create(authority *models.ElectionAuthority) error
    GetByID(id uint) (*models.ElectionAuthority, error)
    GetByElection(electionID uuid.UUID) ([]models.ElectionAuthority, error)
    GetByEmail(email string) (*models.ElectionAuthority, error)
    Update(authority *models.ElectionAuthority) error
    Delete(id uint) error
}

type authoritiesRepository struct {
    db *gorm.DB
}

func NewAuthoritiesRepository(db *gorm.DB) AuthoritiesRepository {
    return &authoritiesRepository{db: db}
}

func (r *authoritiesRepository) Create(authority *models.ElectionAuthority) error {
    return r.db.Create(authority).Error
}

func (r *authoritiesRepository) GetByID(id uint) (*models.ElectionAuthority, error) {
    var authority models.ElectionAuthority
    err := r.db.First(&authority, id).Error
    return &authority, err
}

func (r *authoritiesRepository) GetByElection(electionID uuid.UUID) ([]models.ElectionAuthority, error) {
    var authorities []models.ElectionAuthority
    err := r.db.Where("election = ?", electionID).Find(&authorities).Error
    return authorities, err
}

func (r *authoritiesRepository) GetByEmail(email string) (*models.ElectionAuthority, error) {
    var authority models.ElectionAuthority
    err := r.db.Where("email = ?", email).First(&authority).Error
    return &authority, err
}

func (r *authoritiesRepository) Update(authority *models.ElectionAuthority) error {
    return r.db.Save(authority).Error
}

func (r *authoritiesRepository) Delete(id uint) error {
    return r.db.Delete(&models.ElectionAuthority{}, id).Error
}