package models


import (
	"net"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Ballot struct {
	gorm.Model

	Uuid                    string         `gorm:"type:uuid;default:gen_random_uuid();unique;index" json:"Uuid"`
	Vote                    datatypes.JSON `gorm:"type:jsonb"`
	VotingDeviceFingerprint string         `gorm:"size:225"`
	IPAddress               net.IP         `gorm:"type:inet"`

	StatusID 				uint 			`gorm:"not null" json:"status_id"`
	Status 					StatusBallot 	`gorm:"foreignKey:StatusID;references:ID" json:"-"`
}
