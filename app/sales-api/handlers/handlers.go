package handlers

import (
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) *httptreemux.ContextMux {
	tm := httptreemux.NewContextMux()

	check := check{
		log: log,
	}

	tm.Handle(http.MethodGet, "/test", check.readiness)

	return tm
}
