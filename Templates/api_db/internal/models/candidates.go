package models

import (
    "time"
    "github.com/google/uuid"
)

type Candidate struct {
    ID              int64      `gorm:"primary_key"`
    Elections       uuid.UUID
    Name            string     `gorm:"size:100"`
    Description     string
    PhotoURL        string     `gorm:"size:500;unique"`
    CandidateOrder  int
    CreatedAt       time.Time
    UpdateAt        time.Time
    DeleteAt        *time.Time
}