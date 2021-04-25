CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   name VARCHAR (50) UNIQUE NOT NULL,
   email VARCHAR (255) UNIQUE NOT NULL,
   password_hash TEXT NOT NULL,
   remember_hash TEXT
);