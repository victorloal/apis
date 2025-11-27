package models

import (
    "time"
    "github.com/google/uuid"
)

type Voter struct {
    ID               uint       `gorm:"primary_key"`
    Elections        *uuid.UUID
    Token            string     `gorm:"size:200"`
    VoteStatus       bool       `gorm:"default:false"`
    VerificationHash []byte     `gorm:"type:bytea"`
    IsActive         *bool
    UpdatedAt        time.Time
    CreatedAt        time.Time
    DeleteAt         *time.Time
    
    Ballots []Ballot `gorm:"foreignKey:Voter"`
}