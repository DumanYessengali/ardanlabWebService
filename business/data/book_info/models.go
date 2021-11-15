package book_info

import "time"

type Info struct {
	ID          string    `db:"book_info_id" json:"book_info_id"`
	Year        string    `db:"year" json:"year"`
	Price       string    `db:"price" json:"price"`
	Quantity    int       `db:"quantity" json:"quantity"`
	BookID      string    `db:"book_id" json:"book_id"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewBookInfo struct {
	Year     string `json:"year" validate:"required"`
	Price    string `json:"price" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
	BookID   string `json:"book_id" validate:"required"`
}

type UpdateBookInfo struct {
	Year     *string `json:"year"`
	Price    *string `json:"price"`
	Quantity *int    `json:"quantity"`
	BookID   *string `json:"book_id"`
}
