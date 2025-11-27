package models

import (
    "time"
    "github.com/google/uuid"
)

type ElectionAuditConfig struct {
    ID                  uint      `gorm:"primary_key"`
    Election            uuid.UUID `gorm:"unique"`
    EnableBallotAudit   bool      `gorm:"default:true"`
    EnableAccessLogs    bool      `gorm:"default:true"`
    UpdatedAt           time.Time
    CreatedAt           time.Time
}