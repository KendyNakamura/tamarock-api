CREATE TABLE IF NOT EXISTS article_artist_infos(
	id serial PRIMARY KEY,
	article_id int,
	artist_id int,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)