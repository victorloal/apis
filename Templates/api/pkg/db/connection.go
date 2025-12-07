package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DNS = "host=localhost user=admin password=admin123 dbname=votaciones_db port=5432 sslmode=disable"
var DB *gorm.DB

func DBConetion() {

	var err error
	DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 👈 Tablas en singular
		TablePrefix:   "",   // Sin prefijo
	}})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database connected")
	}

}
