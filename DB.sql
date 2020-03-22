CREATE TABLE IF NOT EXISTS news (
	id serial,
	author text not null,
	body text not null,
	created timestamp NOT NULL,
	PRIMARY KEY (id)
);
