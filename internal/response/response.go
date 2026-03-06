package response

import (
	"encoding/json"
	"net/http"
)

// JSON es el helper universal para respuestas exitosas
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Envolvemos el resultado en una llave "data"
	response := map[string]interface{}{
		"data": data,
	}

	json.NewEncoder(w).Encode(response)
}

// Error es el helper universal para fallos
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"error": message,
	}

	json.NewEncoder(w).Encode(response)
}
