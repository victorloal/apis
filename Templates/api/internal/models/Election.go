package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Election struct {
	gorm.Model

	Uuid        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();unique_index" json:"uuid"`
	Name        string    `gorm:"not null;unique" json:"name"`
	Description string    `gorm:"text" json:"description"`
	StartDate   time.Time    `gorm:"type:timestamp;not null" json:"start_date"`
	EndDate     time.Time    `gorm:"type:timestamp;not null" json:"end_date"`
	Encrypted   bool      `gorm:"default:false" json:"encrypted"`
	Anonymous   bool      `gorm:"default:false" json:"anonumus"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`

	// Foreign keys
	StatusID         uint           `gorm:"not null" json:"status_id"`
	Status           Status         `gorm:"foreignKey:StatusID;references:ID" json:"-"`
	AuditConfigID    uint           `gorm:"default:null" json:"audit_config_id"`
	AuditConfig      AuditConfig    `gorm:"foreignKey:AuditConfigID;references:ID" json:"-"`
	HomomorphicKeyID uint           `gorm:"default:null" json:"homomorphic_key_id"`
	HomomorphicKey   HomomorphicKey `gorm:"foreignKey:HomomorphicKeyID;references:ID" json:"-"`
	TallyResultID    uint           `gorm:"default:null" json:"tally_result_id"`
	TallyResult      TallyResult    `gorm:"foreignKey:TallyResultID;references:ID" json:"-"`

	// inverted foreign keys
	ElectionAuthorities []ElectionAuthority `gorm:"foreignKey:ElectionID" json:"election_authorities"`
	Candidates          []Candidate         `gorm:"foreignKey:ElectionID" json:"candidates"`
	Voters              []Voter             `gorm:"foreignKey:ElectionID" json:"voters"`
}
