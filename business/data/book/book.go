package book

import (
	"context"
	"database/sql"
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/foundation/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/DumanYessengali/ardanlabWebService/business/errors"
	errs "github.com/pkg/errors"
	"log"
)

type Book struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) Book {
	return Book{
		log: log,
		db:  db,
	}
}

func (b Book) Create(ctx context.Context, traceID string, nb NewBook, now time.Time) (Info, error) {
	book := Info{
		ID:          uuid.New().String(),
		Title:       nb.Title,
		Description: nb.Description,
		UserID:      nb.UserID,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO books (book_id, title, description, user_id, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4, $5, $6)`

	b.log.Printf("%s : %s query : %s", traceID, "book.Create",
		database.Log(q, book.ID, book.Title, book.Description, book.UserID, book.DateCreated, book.DateUpdated),
	)

	if _, err := b.db.ExecContext(ctx, q, book.ID, book.Title, book.Description, book.UserID, book.DateCreated, book.DateUpdated); err != nil {
		return Info{}, errs.Wrap(err, "inserting book")
	}
	return book, nil
}

func (b Book) Update(ctx context.Context, traceID string, claims auth.Claims, userID, bookID string, bu UpdateBook, now time.Time) error {
	book, err := b.QueryByID(ctx, traceID, claims, userID, bookID)
	if err != nil {
		return err
	}

	if bu.Title != nil {
		book.Title = *bu.Title
	}
	if bu.Description != nil {
		book.Description = *bu.Description
	}
	book.DateUpdated = now

	const q = `
	UPDATE
		books
	SET 
		"title" = $3,
		"description" = $4,
		"date_updated" = $5
	WHERE
		book_id = $1 AND user_id = $2`

	b.log.Printf("%s: %s: %s", traceID, "book.Update",
		database.Log(q, book.ID, book.Title, book.Description, book.UserID, book.DateCreated, book.DateUpdated),
	)

	if _, err = b.db.ExecContext(ctx, q, book.ID, book.UserID, book.Title, book.Description, book.DateUpdated); err != nil {
		return errs.Wrap(err, "updating book")
	}

	return nil
}

func (b Book) Delete(ctx context.Context, traceID string, bookID string) error {
	if _, err := uuid.Parse(bookID); err != nil {
		return errors.ErrInvalidID
	}
	const q = `DELETE FROM books where book_id = $1`

	b.log.Printf("%s : %s query : %s", traceID, "book.Delete",
		database.Log(q, bookID),
	)

	if _, err := b.db.ExecContext(ctx, q, bookID); err != nil {
		return errs.Wrapf(err, "deleting book %s", bookID)
	}
	return nil
}

func (b Book) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT *FROM books ORDER BY book_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	b.log.Printf("%s : %s query : %s", traceID, "book.Query", database.Log(q, offset, rowsPerPage))

	book := []Info{}

	if err := b.db.SelectContext(ctx, &book, q, offset, rowsPerPage); err != nil {
		return nil, errs.Wrap(err, "selecting book")
	}

	return book, nil
}

func (b Book) QueryByID(ctx context.Context, traceID string, claims auth.Claims, userID, bookID string) (Info, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return Info{}, errors.ErrInvalidID
	}
	if _, err := uuid.Parse(bookID); err != nil {
		return Info{}, errors.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return Info{}, errors.ErrForbidden
	}

	const q = `SELECT * FROM books WHERE user_id = $1 AND book_id = $2`

	b.log.Printf("%s : %s query : %s", traceID, "book.QueryByID",
		database.Log(q, userID),
	)

	var book Info
	if err := b.db.GetContext(ctx, &book, q, userID, bookID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, errors.ErrNotFound
		}
		return Info{}, errs.Wrapf(err, "selecting user %q and book %q", userID, bookID)
	}

	return book, nil
}
