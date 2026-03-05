package main

import (
	"log"
	"net/http"
	"github.com/HaroldVelez13/gohar/internal/handlers"
)

func main() {
	// 1. Instanciamos el handler (nuestro "mini-servicio")
	userHandler := handlers.NewUserHandler()

	// 2. Registramos la ruta
	// En Go, si un struct tiene el método ServeHTTP, puedes pasarlo directo
	http.Handle("/users/", userHandler)

	log.Println("Servidor corriendo en :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}