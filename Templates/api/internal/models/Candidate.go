package models

import (

	"gorm.io/gorm"
)

type Candidate struct {
	gorm.Model

	Name           string `gorm:"not null;unique_index;unique" json:"name"`
	Description    string `gorm:"text" json:"description"`
	PhotoURL       string `json:"photo_url"`
	CandidateOrder int    `gorm:"not null;unique_index" json:"candidate_order"`

	// Foreign keys
	ElectionID uint           `gorm:"not null" json:"election_id"`
	Election   Election `gorm:"foreignKey:ElectionID;references:ID" json:"-"`
}
