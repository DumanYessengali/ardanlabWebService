package handlers

import (
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
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

	//app.Handle(http.MethodGet, "/readiness", cg.readiness, mid.Authenticate(a), mid.Authorized(log, auth.RoleAdmin))
	//app.Handle(http.MethodGet, "/liveness", cg.readiness, mid.Authenticate(a), mid.Authorized(log, auth.RoleAdmin))
	app.Handle(http.MethodGet, "/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/liveness", cg.liveness)
	return app
}
