package models

import (
    "time"
    "github.com/google/uuid"
    "net"
)

type AuditLog struct {
    ID          uint      `gorm:"primary_key"`
    Election    uuid.UUID
    Action      string    `gorm:"size:100"`
    UserType    string    `gorm:"size:50"` // voter, authority, system
    UserID      string    `gorm:"size:100"`
    IPAddress   net.IP    `gorm:"type:inet"`
    UserAgent   string
    Details     JSONB     `gorm:"type:jsonb"`
    CreatedAt   time.Time
}