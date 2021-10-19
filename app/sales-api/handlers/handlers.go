package handlers

import (
	"encoding/json"
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger) *httptreemux.ContextMux {
	tm := httptreemux.NewContextMux()

	h := func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status string
		}{
			Status: "OK",
		}
		json.NewEncoder(w).Encode(status)
	}

	tm.Handle(http.MethodGet, "/test", h)

	return tm
}
