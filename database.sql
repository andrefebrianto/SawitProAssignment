/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL NOT NULL PRIMARY KEY,
	"name" VARCHAR(60) NOT NULL,
  phone VARCHAR(15) NOT NULL UNIQUE,
  "password" VARCHAR(100) NOT NULL,
  login_count BIGINT NOT NULL DEFAULT 0
);


INSERT INTO users ("name", phone, "password") VALUES ('John Doe', '+6281234567890', '$2y$12$wuEg1M4Px4TwENtdJdNx/exizdeEHB1HsiS9SNJ4E0b8msYdJ4HOe');
INSERT INTO users ("name", phone, "password") VALUES ('Mark Levinson', '+6280987654321', '$2y$12$zMdiXZxEF/Two9xcea3OQ.vaLAwmAI9s.LD0TI6gQGU4s/isGdYSy');
