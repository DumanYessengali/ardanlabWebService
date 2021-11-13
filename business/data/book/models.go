package book

import "time"

type Info struct {
	ID          string    `db:"book_id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	UserID      string    `db:"user_id" json:"user_id"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewBook struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
}

type UpdateBook struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
