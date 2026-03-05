package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB ahora recibe la cadena de conexión como parámetro.
// Esto permite que el main decida de dónde viene (env, flag, etc).
func ConnectDB(connStr string) (*pgxpool.Pool, error) {
	// 1. Parsear la configuración de la URL
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("error al parsear config de DB: %v", err)
	}

	// Tip de Senior: Configura límites para el pool
	config.MaxConns = 10                     // Máximo de conexiones abiertas
	config.MinConns = 2                      // Mínimo de conexiones inactivas
	config.MaxConnIdleTime = 5 * time.Minute // Tiempo antes de cerrar una conexión inactiva

	// 2. Crear el pool de conexiones
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error creando el pool: %v", err)
	}

	// 3. Verificar la conexión (Ping) con un timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("no se pudo conectar a la DB: %v", err)
	}

	return pool, nil
}
