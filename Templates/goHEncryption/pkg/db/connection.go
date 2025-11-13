package db

import (
	"fmt"
	"goHEncryption/pkg/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	cfg := config.Load()
	Db := cfg.Db

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Db.Host, Db.User, Db.Password, Db.Dbname, Db.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL using GORM")

	// Migrar las tablas automáticamente (solo en desarrollo)
	err = DB.AutoMigrate(&CrytografiaHomorfica{}, &Publico{}, &AutoridadElectoral{})
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
