package models

import (
	"log"

	"gorm.io/gorm"
)

type StatusBallot struct {
	gorm.Model

	Name        string `gorm:"size:20;unique_index;unique;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

var defaultBallotStatuses = []StatusBallot{
	{Name: "Submitted", Description: "Voto resivido"},
	{Name: "Recorded", Description: "Voto registrado"},
	{Name: "Verified", Description: "Voto verificado"},
	{Name: "Invalid", Description: "Voto inválido"},
	{Name: "Tallied", Description: "Voto contabilizado"},
}

func CreateStatusBallot(db *gorm.DB) {
	for _, status := range defaultBallotStatuses {
		var existing StatusBallot
		// Evitar duplicados
		if err := db.Where("name = ?", status.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&status).Error; err != nil {
					log.Printf("❌ Error creando status '%s': %v", status.Name, err)
				} else {
					log.Printf("✅ Status '%s' creado correctamente", status.Name)
				}
			}
		}
	}
}