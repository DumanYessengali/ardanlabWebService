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
		Version:     2.1,
		Description: "Alter table books with user column",
		Script: `
ALTER TABLE books
	ADD COLUMN user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
`,
	},
}
