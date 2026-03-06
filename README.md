# 🚀 GoHar - User Management Microservice (Enterprise Ready)

Un microservicio de gestión de usuarios desarrollado con **Go 1.24**, diseñado bajo los principios de **Clean Architecture** y **SOLID**. Esta API no es solo un CRUD; es una base profesional con observabilidad, persistencia optimizada y despliegue automatizado.



---

## 🏗️ 1. Arquitectura y Diseño

El proyecto utiliza una estructura de **Separación de Preocupaciones**, permitiendo que la lógica de negocio sea independiente de la base de datos y el servidor web.

* **cmd/api/**: Punto de entrada. Configura e inyecta las dependencias.
* **internal/models/**: Definición de entidades y reglas de validación.
* **internal/handlers/**: Capa de transporte. Maneja peticiones HTTP y respuestas JSON.
* **internal/storage/**: Capa de persistencia. Gestión de Pool de conexiones con PostgreSQL.
* **internal/middleware/**: Capa transversal. Logging, seguridad y recuperación de errores.
* **internal/response/**: Estandarización de la comunicación con el cliente.



---

## 🛠️ 2. Stack Tecnológico

* **Runtime:** Go 1.24+
* **Router:** [Chi v5](https://github.com/go-chi/chi) (Velocidad y compatibilidad total con net/http).
* **Database:** PostgreSQL 15+ (Relacional, robusta y escalable).
* **Driver:** [pgx v5](https://github.com/jackc/pgx) (El estándar de alto rendimiento para Postgres en Go).
* **Validator:** [Go-Playground/Validator](https://github.com/go-playground/validator) (Validación de datos mediante Tags).
* **DevOps:** Docker & Docker Compose (Entornos aislados y reproducibles).

---

## ✨ 3. Características  Implementadas

| Característica | Descripción |
| :--- | :--- |
| **Paginación** | Control de carga mediante `limit` y `offset` con metadatos de página. |
| **Graceful Shutdown** | El servidor cierra conexiones activas antes de apagarse. |
| **Connection Pooling** | Límite de conexiones a DB para evitar saturación de recursos. |
| **Logging Middleware** | Trazabilidad total: `[MÉTODO] [RUTA] [STATUS] [DURACIÓN]`. |
| **JSON Error Handling** | Respuestas consistentes en caso de fallos (404, 400, 500, 409). |
| **Environment Config** | Gestión centralizada de secretos mediante archivos `.env`. |



---

## 🚀 4. Guía de Inicio Rápido

### A. Preparación del Entorno
Clona el repositorio y crea tu archivo de configuración:
```bash
git clone [https://github.com/HaroldVelez13/gohar.git](https://github.com/HaroldVelez13/gohar.git)
cd gohar
touch .env

### B. Configuración del Entorno (.env)
Crea un archivo llamado `.env` en la raíz del proyecto. Este archivo contiene las credenciales y variables que el microservicio necesita para operar en diferentes ambientes:

```env
# Puerto donde correrá la API
PORT=8080

# URL de conexión a PostgreSQL (formato: postgres://usuario:password@host:puerto/nombre_db)
# Si usas Docker, el host debe ser el nombre del servicio definido en docker-compose (ej: 'db')
DB_URL=postgres://user:password@db:5432/micro_db?sslmode=disable

## 🏗️ 5. Despliegue con Docker (Recomendado)

Para garantizar que el microservicio funcione exactamente igual en cualquier máquina, utilizamos **Docker Multi-stage builds** y **Docker Compose**. Esto permite compilar el binario en un entorno controlado y ejecutarlo en una imagen minimalista de Linux (Alpine).

### Pasos para levantar el entorno:

1.  **Compilar y levantar:**
    ```bash
    docker-compose up --build
    ```
2.  **Verificación de salud:**
    La API estará disponible en `http://localhost:8080`.
    Puedes probar la conexión con: `curl http://localhost:8080/users`



---

## 📖 6. Documentación de la API (Endpoints)

### Gestión de Usuarios
Todos los endpoints devuelven un objeto JSON consistente. Si la petición es exitosa, los datos reales vendrán dentro de la llave `"data"`.

| Método | Endpoint | Parámetros Query | Descripción |
| :--- | :--- | :--- | :--- |
| **GET** | `/users` | `page`, `limit` | Obtiene lista de usuarios paginada. |
| **POST** | `/users` | - | Crea un nuevo usuario. |
| **GET** | `/users/{id}` | - | Obtiene el detalle de un usuario por ID. |
| **PUT** | `/users/{id}` | - | Actualiza un usuario existente. |
| **DELETE** | `/users/{id}`| - | Elimina un usuario del sistema. |



### Ejemplos de Respuesta

**Éxito (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "email": "harold.velez@example.com",
      "name": "Harold Velez",
      "age": 32
    }
  ],
  "total": 45,
  "page": 1,
  "last_page": 5
}

**Error de Validación (400 Bad Request):**
```json
{
  "error": "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}

## 🛡️ 7. Seguridad y Resiliencia 

Este microservicio implementa patrones de diseño avanzados para garantizar la estabilidad y protección de los datos en entornos de alta demanda:

* **Prevención de SQL Injection:** No se concatenan strings en las consultas SQL. Se utiliza el motor de parámetros nativo de `pgx` ($1, $2, etc.) para sanear todas las entradas del usuario automáticamente.
* **Panic Recovery:** Implementación de middleware que intercepta errores fatales (como punteros nulos). Esto evita que el contenedor de Go se detenga inesperadamente, respondiendo siempre con un error 500 controlado.
* **Context Control & Timeouts:** Todas las operaciones de base de datos incluyen un `context.Context` con un timeout estricto. Si la base de datos se bloquea o tarda demasiado, el microservicio cancela la operación para liberar recursos y mantenerse disponible.
* **Graceful Shutdown:** El servidor escucha señales del sistema (`SIGINT`, `SIGTERM`). Antes de apagarse, deja de recibir tráfico nuevo, completa las peticiones que ya están en curso y cierra el pool de conexiones a PostgreSQL de forma segura sin corromper datos.



---

## 👨‍💻 Autor
**Harold Velez** - [GitHub](https://github.com/HaroldVelez13) | [LinkedIn](https://linkedin.com/in/tu-perfil)

---
*Este proyecto fue construido siguiendo los estándares de la comunidad de Go en 2026.*