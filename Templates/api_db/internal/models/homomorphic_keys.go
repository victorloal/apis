package models

import (
    "time"
    "github.com/google/uuid"
)

type HomomorphicKey struct {
    ID        uint      `gorm:"primary_key"`
    Elections uuid.UUID
    PKey      []byte    `gorm:"type:bytea;column:p_key"`
    Params    JSONB     `gorm:"type:jsonb"`
    UpdatedAt time.Time
    CreatedAt time.Time
    DeleteAt  *time.Time
}

type JSONB map[string]interface{}