package main

var DBSchema = `
CREATE TABLE IF NOT EXISTS users (
	user_id serial UNIQUE PRIMARY KEY,
	login TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS breeds (
	breed_id serial UNIQUE PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS farms (
	farm_id serial UNIQUE PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT UNIQUE NOT NULL,
	user_id INTEGER NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	CONSTRAINT fk_user_farm
		FOREIGN KEY(user_id) 
		REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS cows (
	cow_id serial UNIQUE PRIMARY KEY,
	name TEXT NOT NULL,
	breed_id INTEGER NOT NULL,
	farm_id INTEGER NOT NULL,
	bolus_sn INTEGER UNIQUE NOT NULL,
	date_of_born DATE NOT NULL,
	added_at TIMESTAMP with time zone NOT NULL,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	CONSTRAINT fk_cow_breed
		FOREIGN KEY(breed_id) 
		REFERENCES breeds(breed_id),
	CONSTRAINT fk_cow_farm
		FOREIGN KEY(farm_id) 
		REFERENCES farms(farm_id)
);

CREATE TABLE IF NOT EXISTS health (
	cow_id INTEGER UNIQUE PRIMARY KEY,
	estrus BOOLEAN DEFAULT FALSE,
	ill TEXT DEFAULT '',
	updated_at TIMESTAMP with time zone DEFAULT NOW(),
	CONSTRAINT fk_health_cow
		FOREIGN KEY(cow_id) 
		REFERENCES cows(cow_id)
);

CREATE TABLE IF NOT EXISTS monitoring_data(
	md_id serial UNIQUE PRIMARY KEY,
	cow_id INTEGER NOT NULL,
	added_at TIMESTAMP with time zone,
	ph FLOAT,
	temperature FLOAT,
	movement FLOAT,
	charge FLOAT,
	CONSTRAINT fk_md_cow
		FOREIGN KEY(cow_id) 
		REFERENCES cows(cow_id)
);
`
