package main

import (
	"encoding/json"
	"net/http"
)

// User: Primera letra mayúscula = Público (Exportado para el serializador JSON)
// Los `struct tags` (entre comillas invertidas) definen el nombre en el JSON
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Age  int8`json:"age"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Creamos el dato (Instanciamos el struct)
	u := User{ID: 1, Email: "dev@example.com", Name:"Developer", Age:37}

	// 2. Definimos el Header (Como en Express: res.set)
	w.Header().Set("Content-Type", "application/json")

	// 3. Serializamos y enviamos (El "Stream" directo al ResponseWriter)
	json.NewEncoder(w).Encode(u)
}

func main() {
	// Definimos la ruta y el handler
	http.HandleFunc("/user", userHandler)

	// Levantamos el servidor (Bloqueante, como app.listen)
	http.ListenAndServe(":8080", nil)
}