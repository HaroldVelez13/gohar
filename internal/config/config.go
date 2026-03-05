package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBURL string
}

func LoadConfig() *Config {
	// Intentamos cargar el .env (en producción se usan variables del sistema)
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró archivo .env, usando variables de entorno del sistema")
	}

	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBURL: getEnv("DB_URL", ""),
	}
}

// Función helper para dar valores por defecto
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
