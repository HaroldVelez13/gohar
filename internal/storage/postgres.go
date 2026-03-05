package storage

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB inicializa un pool de conexiones (mejor que una sola conexión)
func ConnectDB() (*pgxpool.Pool, error) {
	// Obtenemos la URL de la DB del entorno (definida en docker-compose)
	connStr := os.Getenv("DB_URL")
	
	// pgxpool maneja automáticamente la reconexión y el límite de conexiones
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Verificar conexión
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a la DB: %v", err)
	}

	return pool, nil
}