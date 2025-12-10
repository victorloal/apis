package models

import (
	"log"

	"gorm.io/gorm"
)

type Status struct {
	gorm.Model

	Name        string `gorm:"size:20;unique_index;unique;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

var defaultStatuses = []Status{
	{Name: "Draft", Description: "Votación en borrador"},
	{Name: "Scheduled", Description: "Votación programada para una fecha futura"},
	{Name: "Open", Description: "Votación abierta para recibir votos"},
	{Name: "Closed", Description: "Votación cerrada. Ya no se reciben votos"},
	{Name: "Counting", Description: "Conteo o verificación de los votos"},
	{Name: "Completed", Description: "Resultados finalizados y publicados"},
	{Name: "Canceled", Description: "La votación ha sido cancelada"},
}

func CreateStatus(db *gorm.DB) {
	for _, status := range defaultStatuses {
		var existing Status
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