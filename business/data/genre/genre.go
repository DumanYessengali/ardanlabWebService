package genre

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

type Genre struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) Genre {
	return Genre{
		log: log,
		db:  db,
	}
}

func (g Genre) Create(ctx context.Context, traceID string, nb NewGenre, now time.Time) (Info, error) {
	genre := Info{
		ID:          uuid.New().String(),
		Name:        nb.Name,
		BookID:      nb.BookID,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO genres (genre_id, name, book_id, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4, $5)`

	g.log.Printf("%s : %s query : %s", traceID, "genre.Create",
		database.Log(q, genre.ID, genre.Name, genre.BookID, genre.DateCreated, genre.DateUpdated),
	)

	if _, err := g.db.ExecContext(ctx, q, genre.ID, genre.Name, genre.BookID, genre.DateCreated, genre.DateUpdated); err != nil {
		return Info{}, errs.Wrap(err, "inserting genre")
	}
	return genre, nil
}

func (b Genre) Update(ctx context.Context, traceID string, claims auth.Claims, genreID string, bu UpdateGenre, now time.Time) error {
	genre, err := b.QueryByID(ctx, traceID, claims, genreID)
	if err != nil {
		return err
	}

	if bu.Name != nil {
		genre.Name = *bu.Name
	}
	if bu.BookID != nil {
		genre.BookID = *bu.BookID
	}
	genre.DateUpdated = now

	const q = `
	UPDATE
		genres
	SET 
		"name" = $2,
		"book_id" = $3,
		"date_updated" = $4
	WHERE
		book_id = $1`

	b.log.Printf("%s: %s: %s", traceID, "genre.Update",
		database.Log(q, genre.ID, genre.Name, genre.BookID, genre.DateCreated, genre.DateUpdated),
	)

	if _, err = b.db.ExecContext(ctx, q, genre.ID, genre.Name, genre.BookID, genre.DateUpdated); err != nil {
		return errs.Wrap(err, "updating genre")
	}

	return nil
}

func (b Genre) Delete(ctx context.Context, traceID string, genreID string) error {
	if _, err := uuid.Parse(genreID); err != nil {
		return errors.ErrInvalidID
	}
	const q = `DELETE FROM genres where genre_id = $1`

	b.log.Printf("%s : %s query : %s", traceID, "genre.Delete",
		database.Log(q, genreID),
	)

	if _, err := b.db.ExecContext(ctx, q, genreID); err != nil {
		return errs.Wrapf(err, "deleting genre %s", genreID)
	}
	return nil
}

func (b Genre) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT *FROM genres ORDER BY genre_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	b.log.Printf("%s : %s query : %s", traceID, "genre.Query", database.Log(q, offset, rowsPerPage))

	genre := []Info{}

	if err := b.db.SelectContext(ctx, &genre, q, offset, rowsPerPage); err != nil {
		return nil, errs.Wrap(err, "selecting genre")
	}

	return genre, nil
}

func (b Genre) QueryByID(ctx context.Context, traceID string, claims auth.Claims, genreID string) (Info, error) {
	if _, err := uuid.Parse(genreID); err != nil {
		return Info{}, errors.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, errors.ErrForbidden
	}

	const q = `SELECT * FROM genres WHERE genre_id = $1`

	b.log.Printf("%s : %s query : %s", traceID, "genre.QueryByID",
		database.Log(q, genreID),
	)

	var genre Info
	if err := b.db.GetContext(ctx, &genre, q, genreID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, errors.ErrNotFound
		}
		return Info{}, errs.Wrapf(err, "selecting genre %q", genreID)
	}

	return genre, nil
}
