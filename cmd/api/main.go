package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
}

// Simulamos nuestra DB con un Mutex para evitar Race Conditions
var (
	users = []User{
		{ID: 1, Email: "admin@test.com", Name: "Admin", Age: 34},
	}
	mu sync.Mutex 
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// READ: Retornar todos los usuarios
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		// CREATE: Leer el body y agregar
		var newUser User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		mu.Lock() // Bloqueamos para escribir seguro
		newUser.ID = len(users) + 1
		users = append(users, newUser)
		mu.Unlock() // Liberamos

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)

	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/users", handleUsers)
	
	// El servidor ahora corre en el 8080 de tu contenedor
	http.ListenAndServe(":8080", nil)
}