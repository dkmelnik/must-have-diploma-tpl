CREATE TABLE IF NOT EXISTS users (
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     login VARCHAR(255) NOT NULL UNIQUE,
     password VARCHAR(255),
     created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON users (login);