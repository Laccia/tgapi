CREATE TABLE IF NOT EXISTS tghistory (
		id SERIAL PRIMARY KEY,
		msg TEXT NOT NULL,
		msg_id int NOT NULL,
        user_id BIGSERIAL,
		msg_date TIMESTAMP NOT NULL,
		chat_id int NOT NULL);
