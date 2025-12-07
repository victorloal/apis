package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TallyResult struct {
	gorm.Model

	TotalVotes int            `json:"total_votes"`
	Results    datatypes.JSON `gorm:"type:jsonb" json:"results"`
	Skey       []byte         `gorm:"type:bytea;column:s_key" json:"s_key"`
	ComputedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"computed_at"`
}
