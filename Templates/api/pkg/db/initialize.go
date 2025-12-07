package db

import (
	"api/internal/models"
	"log"

	"gorm.io/gorm"
)

func Initialize(db *gorm.DB)  {
	log.Println("🔄 Iniciando base de datos...")

	models.CreateStatusBallot(db)
	models.CreateStatus(db)


	log.Println("✅ Base de datos inicializada")
}