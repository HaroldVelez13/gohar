package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"github.com/go-chi/chi/v5"
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



func (h *UserHandler) GetAll(w http.ResponseWriter, , r *http.Request) {
	json.NewEncoder(w).Encode(h.users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter,  r *http.Request) {
	idStr := chi.URLParam(r, "id") // ¡Mucho más limpio que Split!
	id, _ := strconv.Atoi(idStr)

	for _, u := range h.users {
		if u.ID == id {
			json.NewEncoder(w).Encode(u)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) update(w http.ResponseWriter,   r *http.Request)  {
	idStr := chi.URLParam(r, "id") // ¡Mucho más limpio que Split!
	id, _ := strconv.Atoi(idStr)

	
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

func (h *UserHandler) deleteUser(w http.ResponseWriter,  r *http.Request)  {
	idStr := chi.URLParam(r, "id") // ¡Mucho más limpio que Split!
	id, _ := strconv.Atoi(idStr)
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