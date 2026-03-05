package main

import (
	"log"
	"net/http"

	"github.com/HaroldVelez13/gohar/internal/config"
	"github.com/HaroldVelez13/gohar/internal/handlers"
	"github.com/HaroldVelez13/gohar/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Cargar Configuración
	cfg := config.LoadConfig()

	// 2. Conectar a DB usando la URL de la config
	db, err := storage.ConnectDB(cfg.DBURL) //
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Inyectar DB al handler
	userH := handlers.NewUserHandler(db)

	r := chi.NewRouter()

	// Middlewares: Chi trae algunos de fábrica muy útiles
	r.Use(middleware.Logger)    // Loguea cada petición en consola
	r.Use(middleware.Recoverer) // Evita que el server muera si hay un panic

	// Definición de rutas tipo Express
	r.Route("/users", func(r chi.Router) {
		r.Get("/", userH.GetAll)
		r.Post("/", userH.Create)
		r.Get("/{id}", userH.GetByID) // Parámetro de URL nombrado
		r.Put("/{id}", userH.Update)
		r.Delete("/{id}", userH.Delete)
	})

	log.Println("Servidor iniciado en :8080")
	http.ListenAndServe(":8080", r)
}
