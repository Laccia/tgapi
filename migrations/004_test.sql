CREATE TABLE IF NOT EXISTS tgusers (
		id SERIAL PRIMARY KEY,
        user_id BIGSERIAL NOT NULL UNIQUE,
        first_name TEXT,
        last_name TEXT,
        username TEXT,
        phone TEXT,
        premium BOOLEAN);
