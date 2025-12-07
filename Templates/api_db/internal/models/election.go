package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Election struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Status      int
    Encrypted   bool
    Name        string
    Description string
    StartDate   time.Time
    EndDate     time.Time
    UpdatedAt   time.Time
    CreatedAt   time.Time
    
    // Relations
    Authorities    []ElectionAuthority `gorm:"foreignKey:Election"`
    HomomorphicKey HomomorphicKey     `gorm:"foreignKey:Elections"`
    Candidates     []Candidate        `gorm:"foreignKey:Elections"`
    Voters         []Voter            `gorm:"foreignKey:Elections"`
    Ballots        []Ballot           `gorm:"foreignKey:Elections"`
    TallyResult    TallyResult        `gorm:"foreignKey:Election"`
    AuditConfig    ElectionAuditConfig `gorm:"foreignKey:Election"`
}

func (e *Election) BeforeCreate(tx *gorm.DB) error {
    if e.ID == uuid.Nil {
        e.ID = uuid.New()
    }
    return nil
}