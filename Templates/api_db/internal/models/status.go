package models

import "time"

type Status struct {
    ID        int       `gorm:"primary_key"`
    Name      string    `gorm:"size:20;unique"`
    UpdatedAt time.Time
    CreatedAt time.Time
    DeleteAt  *time.Time
}