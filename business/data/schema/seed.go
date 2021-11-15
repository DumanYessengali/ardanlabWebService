package schema

import "github.com/jmoiron/sqlx"

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

const seeds = `
INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;

INSERT INTO books (book_id, title, description, user_id, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'three cups of tea', 'desc1', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'it', 'desc2','45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO book_infos (book_info_id, year, price, quantity, book_id, date_created, date_updated) VALUES
	('1150e3bf-19c2-4de2-9c41-eed8e5b23b06', '1942', '123', 12, 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('2826a093-4fe0-4519-bd3e-e43f5d4afede', '1992', '123', 12, '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO genres (genre_id, name, book_id, date_created, date_updated) VALUES
	('5c4ad2f0-b4f8-4cb2-ae5b-15d5ddd45e45', 'horror', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('d9abd828-bcf5-4d54-a7fd-efa666897181', 'comedy', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;`

const deleteAll = `
DELETE FROM genres;
DELETE FROM book_infos;
DELETE FROM books;
DELETE FROM users;`

func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
