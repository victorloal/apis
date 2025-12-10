package models

import (
	"gorm.io/gorm"
)

type AuditConfig struct {
	gorm.Model

	EnableBallotAudit bool `gorm:"default:true" json:"enable_ballot_audit"`
	EnableAccessLogs  bool `gorm:"default:true" json:"enable_access_logs"`
}
