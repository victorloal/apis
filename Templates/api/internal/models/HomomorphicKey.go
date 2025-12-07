package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type HomomorphicKey struct {
	gorm.Model

	PKey   []byte         `gorm:"type:bytea;column:p_key"`
	Params datatypes.JSON `gorm:"type:jsonb"`
}
