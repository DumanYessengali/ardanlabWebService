package genre

import "time"

type Info struct {
	ID          string    `db:"genre_id" json:"genre_id"`
	Name        string    `db:"name" json:"name"`
	BookID      string    `db:"book_id" json:"book_id"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewGenre struct {
	Name   string `json:"name" validate:"required"`
	BookID string `json:"book_id" validate:"required"`
}

type UpdateGenre struct {
	Name   *string `json:"name"`
	BookID *string `json:"book_id"`
}
