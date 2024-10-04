CREATE TABLE IF NOT EXISTS tghistory (
		id SERIAL PRIMARY KEY,
		msg TEXT NOT NULL,
		msg_id int NOT NULL,
		msg_date TIMESTAMP NOT NULL,
		chat_id int NOT NULL);
