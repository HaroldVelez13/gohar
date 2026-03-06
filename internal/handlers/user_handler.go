package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/HaroldVelez13/gohar/internal/models"
	"github.com/HaroldVelez13/gohar/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	db       *pgxpool.Pool
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
	// 1. Obtener parámetros de paginación con valores por defecto
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 5
	} // Máximo 100 por seguridad

	offset := (page - 1) * limit

	// 2. Obtener el total de registros para la metadata
	var total int
	err := h.db.QueryRow(r.Context(), "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al contar registros")
		return
	}

	// 3. Query con LIMIT y OFFSET
	rows, err := h.db.Query(r.Context(),
		"SELECT id, email, name, age FROM users ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener datos")
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Email, &u.Name, &u.Age)
		users = append(users, u)
	}

	// 4. Construir respuesta final
	lastPage := (total + limit - 1) / limit
	response := models.UserPagedResponse{
		Data:     users,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
			response.Error(w, http.StatusNotFound, "Usuario no encontrado")
		} else {
			response.Error(w, http.StatusInternalServerError, err.Error())
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
		response.Error(w, http.StatusBadRequest, "JSON mal formado")
		return
	}

	// 2. Validación de reglas
	if err := h.validate.Struct(u); err != nil {
		// Si hay error, devolvemos un 400 con el detalle
		response.Error(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
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
			response.Error(w, http.StatusConflict, "El correo electrónico ya está registrado")
			return
		}

		// Cualquier otro error de DB sí es un 500
		response.Error(w, http.StatusInternalServerError, "Error inesperado al guardar en la base de datos")
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
		response.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// 1. Validar Formato (Email, Min Length, etc.)
	if err := h.validate.Struct(u); err != nil {
		response.Error(w, http.StatusBadRequest, "Validación fallida: "+err.Error())
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
			response.Error(w, http.StatusConflict, "El email ya está en uso por otro usuario")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error interno")
		return
	}

	if tag.RowsAffected() == 0 {
		response.Error(w, http.StatusNotFound, "Usuario no encontrado")
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
