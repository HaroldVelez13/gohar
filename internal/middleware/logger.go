package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Usamos un 'WrapResponseWriter' para poder leer el Status Code al final,
		// ya que el ResponseWriter estándar de Go no permite leerlo después de escribirlo.
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Ejecutar el siguiente handler en la cadena
		next.ServeHTTP(ww, r)

		// Calcular duración
		duration := time.Since(start)

		// Imprimir log formateado
		log.Printf(
			"[%s] %s | %d | %s",
			r.Method,
			r.URL.Path,
			ww.Status(),
			duration,
		)
	})
}
