CREATE TABLE IF NOT EXISTS documents(
   id serial PRIMARY KEY,
   text VARCHAR (50) UNIQUE NOT NULL,
   user_id serial REFERENCES users(id)
);