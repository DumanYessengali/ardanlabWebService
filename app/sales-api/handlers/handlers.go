package handlers

import (
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/business/data/book"
	"github.com/DumanYessengali/ardanlabWebService/business/data/book_info"
	"github.com/DumanYessengali/ardanlabWebService/business/data/genre"
	"github.com/DumanYessengali/ardanlabWebService/business/data/user"
	"github.com/DumanYessengali/ardanlabWebService/business/mid"
	"github.com/DumanYessengali/ardanlabWebService/foundation/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	cg := CheckGroup{
		build: build,
		db:    db,
	}

	app.Handle(http.MethodGet, "/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/liveness", cg.liveness)

	ug := userGroup{
		user: user.New(log, db),
		auth: a,
	}
	app.Handle(http.MethodGet, "/users/:page/:rows", ug.query, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodGet, "/users/:id", ug.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/users/token/:kid", ug.token)
	app.Handle(http.MethodPost, "/users", ug.create, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodPut, "/users/:id", ug.update, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/users/:id", ug.delete, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))

	bg := bookGroup{
		book: book.New(log, db),
		auth: a,
	}
	app.Handle(http.MethodGet, "/books/:page/:rows", bg.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/books/:book_id", bg.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/books", bg.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/books/:id", bg.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/books/:id", bg.delete, mid.Authenticate(a))

	big := bookInfoGroup{
		bookInfo: book_info.New(log, db),
		auth:     a,
	}
	app.Handle(http.MethodGet, "/books-info/:page/:rows", big.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/books-info/:book_id", big.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/books-info", big.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/books-info/:id", big.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/books-info/:id", big.delete, mid.Authenticate(a))

	gg := genreGroup{
		genre: genre.New(log, db),
		auth:  a,
	}
	app.Handle(http.MethodGet, "/genre/:page/:rows", gg.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/genre/:book_id", gg.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/genre", gg.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/genre/:id", gg.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/genre/:id", gg.delete, mid.Authenticate(a))

	return app
}
