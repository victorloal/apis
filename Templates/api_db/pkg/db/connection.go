package db

import (
    "api_db/pkg/config"
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
        cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return db, nil
}