CREATE TABLE IF NOT EXISTS articles(
	id serial PRIMARY KEY,
	title VARCHAR (255) UNIQUE NOT NULL,
	text TEXT,
	category INT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at DATETIME DEFAULT NULL
)