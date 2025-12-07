package models

import (
    "time"
    "github.com/google/uuid"
    "net"
)

type Ballot struct {
    ID                     string    `gorm:"primary_key;size:200"`
    Elections              uuid.UUID `gorm:"primary_key"`
    Voter                  int       `gorm:"primary_key"`
    Vote                   JSONB     `gorm:"type:jsonb"`
    VotingDeviceFingerprint string   `gorm:"size:225"`
    IPAddress              net.IP    `gorm:"type:inet"`
    UpdatedAt              time.Time
    CreatedAt              time.Time
    DeleteAt               *time.Time
}