package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1.1,
		Description: "Create table users",
		Script: `
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,
	PRIMARY KEY (user_id)
);`,
	},
	{
		Version:     1.2,
		Description: "Create table books",
		Script: `
CREATE TABLE books (
	book_id   	UUID,
	title       TEXT,
	description TEXT,
	user_id     UUID,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,
	PRIMARY KEY (book_id)
);`,
	},
	{
		Version:     1.3,
		Description: "Create table book_info",
		Script: `
CREATE TABLE book_infos (
	book_info_id UUID,
	year         TEXT,
	price 		 TEXT,
	quantity 	 INTEGER,
	book_id      UUID,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,
	PRIMARY KEY (book_info_id)
);`,
	},
	{
		Version:     1.4,
		Description: "Create table genre",
		Script: `
CREATE TABLE genres (
	genre_id UUID,
	name         TEXT,
	book_id      UUID,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,
	PRIMARY KEY (genre_id)
);`,
	},
}
