package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"github.com/HaroldVelez13/gohar/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	db *pgxpool.Pool
	validate *validator.Validate // Instancia del validador
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		db:       db,
		validate: validator.New(), // Inicializamos aquí
	}
}

// GET /users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), "SELECT id, email, name, age FROM users ORDER BY id")
	if err != nil {
		http.Error(w, "Error al consultar usuarios", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []models.User{} // Inicializamos vacío para que el JSON no sea 'null' sino '[]'
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Age); err != nil {
			continue
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GET /users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var u models.User
	err := h.db.QueryRow(r.Context(), 
		"SELECT id, email, name, age FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Email, &u.Name, &u.Age)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// POST /users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u models.User
	
	// 1. Decode del JSON
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "JSON mal formado", http.StatusBadRequest)
		return
	}

	// 2. Validación de reglas
	if err := h.validate.Struct(u); err != nil {
		// Si hay error, devolvemos un 400 con el detalle
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Si todo está bien, procedemos al INSERT
	err := h.db.QueryRow(r.Context(), 
		"INSERT INTO users (email, name, age) VALUES ($1, $2, $3) RETURNING id", 
		u.Email, u.Name, u.Age).Scan(&u.ID)

	if err != nil {
        // 4. Capturar el error de Email Duplicado (Unique Constraint)
        // El código de error estándar de Postgres para unique_violation es 23505
        if strings.Contains(err.Error(), "unique_violation") || strings.Contains(err.Error(), "23505") {
            h.sendError(w, "El correo electrónico ya está registrado", http.StatusConflict)
            return
        }

        // Cualquier otro error de DB sí es un 500
        h.sendError(w, "Error inesperado al guardar en la base de datos", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// PUT /users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.sendError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 1. Validar Formato (Email, Min Length, etc.)
	if err := h.validate.Struct(u); err != nil {
		h.sendError(w, "Validación fallida: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Ejecutar el Update en la DB
	// Usamos el ID del URL ($4) y los datos del Body ($1, $2, $3)
	tag, err := h.db.Exec(r.Context(), 
		"UPDATE users SET email=$1, name=$2, age=$3 WHERE id=$4", 
		u.Email, u.Name, u.Age, id)

	if err != nil {
		// 3. Validar si el email ya existe (Error de Postgres: 23505)
		// En pgx, puedes capturar errores específicos de la DB
		if strings.Contains(err.Error(), "unique_violation") || strings.Contains(err.Error(), "23505") {
			h.sendError(w, "El email ya está en uso por otro usuario", http.StatusConflict)
			return
		}
		h.sendError(w, "Error interno", http.StatusInternalServerError)
		return
	}

	if tag.RowsAffected() == 0 {
		h.sendError(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	u.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// DELETE /users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	tag, err := h.db.Exec(r.Context(), "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tag.RowsAffected() == 0 {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}