package db

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CrytografiaHomorfica struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	T  int       `gorm:"not null"`
	N  int       `gorm:"not null"`

	Publico Publico `gorm:"foreignKey:HE;references:ID"`
}

type Publico struct {
	ID     uuid.UUID       `gorm:"type:uuid;primaryKey"`
	PK     string          `gorm:"type:varchar(255)"`
	HE     uuid.UUID       `gorm:"type:uuid"`
	Params json.RawMessage `gorm:"type:jsonb"`

	Autoridades []AutoridadElectoral `gorm:"foreignKey:HE;references:HE"`
}

type AutoridadElectoral struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string
	Mail     string
	Password string
	IdParty  uuid.UUID `gorm:"type:uuid"`
	HE       uuid.UUID `gorm:"type:uuid"`
}
