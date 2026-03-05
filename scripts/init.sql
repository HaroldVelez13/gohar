-- 1. Crear la tabla con ID autoincrementable
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, -- Esto genera el autoincremento automáticamente
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    age SMALLINT
);

-- 2. Insertar 10 usuarios de prueba
INSERT INTO users (email, name, age) VALUES
('harold.fofo@example.com', 'Harold Fofo', 32),
('laura.remia@example.com', 'Laura Remia', 28),
('carlos.perez@example.com', 'Carlos Perez', 45),
('ana.martinez@example.com', 'Ana Martinez', 22),
('sergio.rojas@example.com', 'Sergio Rojas', 35),
('marta.lopez@example.com', 'Marta Lopez', 29),
('diego.ruiz@example.com', 'Diego Ruiz', 41),
('elena.sanz@example.com', 'Elena Sanz', 26),
('jorge.castro@example.com', 'Jorge Castro', 38),
('sofia.duarte@example.com', 'Sofia Duarte', 31)
ON CONFLICT (email) DO NOTHING; -- Evita errores si reinicias el contenedor