package models

import (
	"gorm.io/gorm"
)

type ElectionAuthority struct {
	gorm.Model

	CC       uint   `gorm:"not null; unique_index" json:"cc"`
	Name     string `gorm:"size:100"`
	Email    string `gorm:"size:255;unique"`
	Password string `gorm:"size:255"`
	Phone    string `json:"phone"`
	Role     string `gorm:"not null" json:"role"`
	SKey     []byte `gorm:"type:bytea;column:s_key"`
	IsActive bool   `gorm:"default:true" json:"is_active"`

	// Foreign keys
	ElectionID uint     `gorm:"not null" json:"election_id"`
	Election   Election `gorm:"foreignKey:ElectionID;references:ID" json:"-"`
}
