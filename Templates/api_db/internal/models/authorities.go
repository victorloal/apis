package models

import (
    "time"
    "github.com/google/uuid"
)

type ElectionAuthority struct {
    ID        uint      `gorm:"primary_key"`
    Election  uuid.UUID
    CC        int
    Name      string    `gorm:"size:100"`
    Email     string    `gorm:"size:255;unique"`
    Password  string    `gorm:"size:255"`
    SKey      []byte    `gorm:"type:bytea;column:s_key"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeleteAt  *time.Time
}