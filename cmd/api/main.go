package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
}

var (
	users = []User{
		{ID: 1, Email: "admin@test.com", Name: "Admin", Age: 30},
	}
	mu sync.Mutex
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Separamos el path para buscar el ID: /users/1 -> ["", "users", "1"]
	parts := strings.Split(r.URL.Path, "/")
	var id int
	if len(parts) > 2 && parts[2] != "" {
		id, _ = strconv.Atoi(parts[2])
	}

	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			// GET BY ID
			for _, u := range users {
				if u.ID == id {
					json.NewEncoder(w).Encode(u)
					return
				}
			}
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// GET ALL
			json.NewEncoder(w).Encode(users)
		}

	case http.MethodPost:
		var newUser User
		json.NewDecoder(r.Body).Decode(&newUser)
		mu.Lock()
		newUser.ID = len(users) + 1
		users = append(users, newUser)
		mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)

	case http.MethodPut:
		if id == 0 {
			http.Error(w, "ID requerido", http.StatusBadRequest)
			return
		}
		var updatedData User
		json.NewDecoder(r.Body).Decode(&updatedData)

		mu.Lock()
		defer mu.Unlock() // Se ejecuta al final de la función automáticamente
		for i, u := range users {
			if u.ID == id {
				updatedData.ID = id // Aseguramos que el ID no cambie
				users[i] = updatedData
				json.NewEncoder(w).Encode(updatedData)
				return
			}
		}
		http.Error(w, "No encontrado", http.StatusNotFound)

	case http.MethodDelete:
		if id == 0 {
			http.Error(w, "ID requerido", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for i, u := range users {
			if u.ID == id {
				// Truco de Go para borrar un elemento de un slice
				users = append(users[:i], users[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "No encontrado", http.StatusNotFound)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Usamos /users/ para que acepte sub-rutas como /users/1
	http.HandleFunc("/users/", handleUsers)
	http.ListenAndServe(":8080", nil)
}