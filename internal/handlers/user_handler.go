package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"github.com/HaroldVelez13/gohar/internal/models" // Importa tu modelo
)

type UserHandler struct {
	mu    sync.Mutex
	users []models.User
}

// NewUserHandler es un "Constructor" (patrón común en Go)
func NewUserHandler() *UserHandler {
	return &UserHandler{
		users: []models.User{
			{ID: 1, Email: "admin@test.com", Name: "Admin", Age: 30},
		},
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	parts := strings.Split(r.URL.Path, "/")
	var id int
	if len(parts) > 2 && parts[2] != "" {
		id, _ = strconv.Atoi(parts[2])
	}

	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, id)
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodPut:
		h.updateUser(w,id,r)
	case http.MethodDelete:
		h.deleteUser(w, id)
	// ... aquí irían Delete y Update siguiendo el mismo patrón
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) getUsers(w http.ResponseWriter, id int) {
	if id > 0 {
		for _, u := range h.users {
			if u.ID == id {
				json.NewEncoder(w).Encode(u)
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(h.users)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	
	newUser.ID = len(h.users) + 1
	h.users = append(h.users, newUser)
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, id int,  r *http.Request)  {
	if id == 0 {
			http.Error(w, "ID requerido", http.StatusBadRequest)
			return
		}
		var updateUser models.User
		json.NewDecoder(r.Body).Decode(&updateUser)

		h.mu.Lock()
		defer h.mu.Unlock() // Se ejecuta al final de la función automáticamente
		for i, u := range h.users {
			if u.ID == id {
				updateUser.ID = id // Aseguramos que el ID no cambie
				h.users[i] = updateUser
				json.NewEncoder(w).Encode(updateUser)
				return
			}
		}
		http.Error(w, "No encontrado", http.StatusNotFound)
	
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, id int) {
	if id == 0 {
			http.Error(w, "ID requerido", http.StatusBadRequest)
			return
		}
		h.mu.Lock()
		defer h.mu.Unlock()
		for i, u := range h.users {
			if u.ID == id {
				// Truco de Go para borrar un elemento de un slice
				h.users = append(h.users[:i], h.users[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "No encontrado", http.StatusNotFound)
	
}