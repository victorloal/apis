package db

import (
	"api/internal/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("🔄 Ejecutando migraciones...")

	// Orden correcto de migración (sin dependencias circulares)
	err := db.AutoMigrate(
		&models.Status{},
		&models.AuditConfig{},
		&models.Election{},
		&models.ElectionAuthority{},
		&models.Candidate{},
		&models.Voter{},
		&models.HomomorphicKey{},
		&models.TallyResult{},
		&models.Ballot{},
		&models.StatusBallot{},
	)

	if err != nil {
		return err
	}


	log.Println("✅ Migraciones completadas")
	return nil
}
