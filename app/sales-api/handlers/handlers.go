package handlers

import (
	"github.com/DumanYessengali/ardanlabWebService/business/mid"
	"github.com/DumanYessengali/ardanlabWebService/foundation/web"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log))

	check := Check{
		log: log,
	}

	app.Handle(http.MethodGet, "/readiness", check.readiness)

	return app
}
