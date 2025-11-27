package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config representa la estructura completa de configuración
type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	Features FeaturesConfig `yaml:"features"`
	He Hencrytor `yaml:"he"`
	Mail Mail `yaml:"mail"`
	Db Db `yaml:"db"`
}

type Hencrytor struct{
	Parties int `yaml:"parties"`
	Threshold int `yaml:"threshold"`
	LogN int `yaml:"logN"`
	LogP int `yaml:"logP"`
	LogQ int `yaml:"logQ"`
	Module int `yaml:"module"`
}

type Db struct{
	Host string `yaml:"DB_HOST"`
	Port string `yaml:"DB_PORT"`
	User string `yaml:"DB_USER"`
	Password string `yaml:"DB_PASSWORD"`
	Dbname string `yaml:"DB_NAME"`
}

type Mail struct{
	Url string `yaml:"url"`
}

// AppConfig configuración de la aplicación
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// FeaturesConfig configuración de features
type FeaturesConfig struct {
	EnableSwagger bool `yaml:"enable_swagger"`
	EnableMetrics bool `yaml:"enable_metrics"`
}

var (
	// Global config instance
	globalConfig *Config
)

// Load carga la configuración desde el archivo YAML
func Load() *Config {
	if globalConfig != nil {
		return globalConfig
	}

	// Buscar el archivo de configuración en diferentes ubicaciones
	configPaths := []string{
		"configs/config.yaml",
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		log.Fatal("No se encontró el archivo de configuración config.yaml")
	}

	// Leer archivo
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error leyendo archivo de configuración: %v", err)
	}

	// Parsear YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Fatalf("Error parseando YAML: %v", err)
	}

	// Sobrescribir con variables de entorno si existen
	config.loadFromEnv()

	globalConfig = config
	log.Printf("Configuración cargada desde: %s", configFile)
	
	return globalConfig
}

// loadFromEnv sobrescribe configuración con variables de entorno
func (c *Config) loadFromEnv() {
	if port := os.Getenv("APP_PORT"); port != "" {
		c.Server.Port = port
	}
	
	if env := os.Getenv("APP_ENV"); env != "" {
		c.App.Environment = env
	}
}

// Get retorna la instancia global de configuración
func Get() *Config {
	if globalConfig == nil {
		return Load()
	}
	return globalConfig
}

// IsDevelopment verifica si estamos en entorno de desarrollo
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction verifica si estamos en entorno de producción
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}