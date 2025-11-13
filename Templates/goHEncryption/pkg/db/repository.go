package db

import (
	"encoding/json"
	"fmt"
	"goHEncryption/pkg/config"

	"github.com/google/uuid"
)

// Guarda una nueva instancia de CrytografiaHomorfica
func SaveCrytografiaHomorfica() error {
	cfg := config.Load()
	he := CrytografiaHomorfica{
		ID: uuid.New(),
		T:  cfg.He.Threshold,
		N:  cfg.He.Parties,
	}
	if err := DB.Create(&he).Error; err != nil {
		return fmt.Errorf("error inserting CrytografiaHomorfica: %w", err)
	}
	return nil
}

// Guarda una nueva fila en la tabla Publico
func SavePublico(pk string, he uuid.UUID, params json.RawMessage) error {
	public := Publico{
		ID:     uuid.New(),
		PK:     pk,
		HE:     he,
		Params: params,
	}
	if err := DB.Create(&public).Error; err != nil {
		return fmt.Errorf("error inserting Publico: %w", err)
	}
	return nil
}

// Guarda una nueva autoridad electoral
func SaveAE(name, mail, password string, he, idParty uuid.UUID) error {
	ae := AutoridadElectoral{
		ID:       uuid.New(),
		Name:     name,
		Mail:     mail,
		Password: password,
		IdParty:  idParty,
		HE:       he,
	}
	if err := DB.Create(&ae).Error; err != nil {
		return fmt.Errorf("error inserting AutoridadElectoral: %w", err)
	}
	return nil
}

// Obtiene un registro de CrytografiaHomorfica
func GetCrytografiaHomorfica() ([]CrytografiaHomorfica, error) {
	var list []CrytografiaHomorfica
	if err := DB.First(&list).Error; err != nil {
		return nil, fmt.Errorf("error getting CrytografiaHomorfica: %w", err)
	}
	return list, nil
}

// Obtiene un registro de Publico según el HE
func GetPublico(he uuid.UUID) (Publico, error) {
	var public Publico
	if err := DB.Where("he = ?", he).First(&public).Error; err != nil {
		return Publico{}, fmt.Errorf("error getting Publico: %w", err)
	}
	return public, nil
}

// obtiene un registro de AutoridadesElectorales 
func GetAE(he uuid.UUID) ([]AutoridadElectoral, error) {
	var list []AutoridadElectoral
	if err := DB.Where("he = ?", he).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("error getting AutoridadesElectorales: %w", err)
	}
	return list, nil
}

// modifica un registro de Publico segun el HE
func UpdatePublico(he uuid.UUID, data map[string]interface{}) error {
	if err := DB.Model(&Publico{}).Where("he = ?", he).Updates(data).Error; err != nil {
		return fmt.Errorf("error updating Publico: %w", err)
	}
	return nil
}

func UpdateAE(he uuid.UUID, data map[string]interface{}) error {
	if err := DB.Model(&AutoridadElectoral{}).Where("he = ?", he).Updates(data).Error; err != nil {
		return fmt.Errorf("error updating AutoridadesElectorales: %w", err)
	}
	return nil
}
