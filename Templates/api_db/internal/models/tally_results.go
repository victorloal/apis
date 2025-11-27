package models

import (
    "time"
    "github.com/google/uuid"
)

type TallyResult struct {
    ID          uint      `gorm:"primary_key"`
    Election    uuid.UUID `gorm:"unique"`
    Results     JSONB     `gorm:"type:jsonb"`
    TotalVotes  int
    ComputedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    ComputedBy  string    `gorm:"size:255"`
    Proof       JSONB     `gorm:"type:jsonb"`
    UpdatedAt   time.Time
    CreatedAt   time.Time
}