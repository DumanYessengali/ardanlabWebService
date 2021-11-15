package book_info

import (
	"context"
	"database/sql"
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/business/errors"
	"github.com/DumanYessengali/ardanlabWebService/foundation/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	errs "github.com/pkg/errors"
	"log"
	"time"
)

type BookInfo struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) BookInfo {
	return BookInfo{
		log: log,
		db:  db,
	}
}

func (b BookInfo) Create(ctx context.Context, traceID string, nb NewBookInfo, now time.Time) (Info, error) {
	book := Info{
		ID:          uuid.New().String(),
		Year:        nb.Year,
		Price:       nb.Price,
		Quantity:    nb.Quantity,
		BookID:      nb.BookID,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO book_infos (book_info_id, year, price, quantity, book_id, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4, $5, $6, $7)`

	b.log.Printf("%s : %s query : %s", traceID, "bookInfo.Create",
		database.Log(q, book.ID, book.Year, book.Price, book.Quantity, book.BookID, book.DateCreated, book.DateUpdated),
	)

	if _, err := b.db.ExecContext(ctx, q, book.ID, book.Year, book.Price, book.Quantity, book.BookID, book.DateCreated, book.DateUpdated); err != nil {
		return Info{}, errs.Wrap(err, "inserting book info")
	}
	return book, nil
}

func (b BookInfo) Update(ctx context.Context, traceID string, claims auth.Claims, bookInfoID string, bu UpdateBookInfo, now time.Time) error {
	book, err := b.QueryByID(ctx, traceID, claims, bookInfoID)
	if err != nil {
		return err
	}

	if bu.Year != nil {
		book.Year = *bu.Year
	}
	if bu.Price != nil {
		book.Price = *bu.Price
	}
	if bu.BookID != nil {
		book.BookID = *bu.BookID
	}
	book.DateUpdated = now

	const q = `
	UPDATE
		book_info
	SET 
		"year" = $2,
		"price" = $3,
	    "quantity" = $4
		"book_id" = $5,
		"date_updated" = $6
	WHERE
		book_id = $1`

	b.log.Printf("%s: %s: %s", traceID, "bookInfo.Update",
		database.Log(q, book.ID, book.Year, book.Price, book.Quantity, book.BookID, book.DateCreated, book.DateUpdated),
	)

	if _, err = b.db.ExecContext(ctx, q, book.ID, book.Year, book.Price, book.Quantity, book.BookID, book.DateUpdated); err != nil {
		return errs.Wrap(err, "updating bookInfo")
	}

	return nil
}

func (b BookInfo) Delete(ctx context.Context, traceID string, bookInfoID string) error {
	if _, err := uuid.Parse(bookInfoID); err != nil {
		return errors.ErrInvalidID
	}
	const q = `DELETE FROM book_infos where book_info_id = $1`

	b.log.Printf("%s : %s query : %s", traceID, "bookInfo.Delete",
		database.Log(q, bookInfoID),
	)

	if _, err := b.db.ExecContext(ctx, q, bookInfoID); err != nil {
		return errs.Wrapf(err, "deleting book info %s", bookInfoID)
	}
	return nil
}

func (b BookInfo) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT *FROM book_infos ORDER BY book_info_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	b.log.Printf("%s : %s query : %s", traceID, "bookInfo.Query", database.Log(q, offset, rowsPerPage))

	book := []Info{}

	if err := b.db.SelectContext(ctx, &book, q, offset, rowsPerPage); err != nil {
		return nil, errs.Wrap(err, "selecting book info")
	}

	return book, nil
}

func (b BookInfo) QueryByID(ctx context.Context, traceID string, claims auth.Claims, bookInfoID string) (Info, error) {
	if _, err := uuid.Parse(bookInfoID); err != nil {
		return Info{}, errors.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, errors.ErrForbidden
	}

	const q = `SELECT * FROM book_infos WHERE book_info_id = $1`

	b.log.Printf("%s : %s query : %s", traceID, "bookInfo.QueryByID",
		database.Log(q, bookInfoID),
	)

	var book Info
	if err := b.db.GetContext(ctx, &book, q, bookInfoID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, errors.ErrNotFound
		}
		return Info{}, errs.Wrapf(err, "selecting book %q", bookInfoID)
	}

	return book, nil
}
