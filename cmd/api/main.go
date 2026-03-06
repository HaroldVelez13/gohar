package main

import (
	"log"
	"net/http"

	"github.com/HaroldVelez13/gohar/internal/config"
	"github.com/HaroldVelez13/gohar/internal/handlers"
	customMW "github.com/HaroldVelez13/gohar/internal/middleware"
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

	r := chi.NewRouter()

	// 2. Inyectar DB al handler

	// Middlewares: Chi trae algunos de fábrica muy útiles
	// --- MIDDLEWARES GLOBALES ---
	r.Use(middleware.Recoverer) // Evita que el servidor muera si hay un panic
	r.Use(customMW.Logger)      // Nuestro nuevo logger (ajusta el import)
	r.Use(middleware.RealIP)    // Obtiene la IP real del cliente

	userH := handlers.NewUserHandler(db)
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
