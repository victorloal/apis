package models

import (
	"gorm.io/gorm"
)

type Voter struct {
	gorm.Model

	Token      string `gorm:"not null;unique_index;unique" json:"token"`
	VoteStatus bool   `gorm:"default:false" json:"vote_status"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`

	// Foreign keys
	ElectionID uint     `gorm:"not null" json:"election_id"`
	Election   Election `gorm:"foreignKey:ElectionID;references:ID" json:"-"`
	BallotID   uint     `gorm:"default:null" json:"ballot_id"`
	Ballot     Ballot   `gorm:"foreignKey:BallotID;references:ID" json:"-"`
}
