package handlers

import (
	"context"
	"errors"
	"github.com/DumanYessengali/ardanlabWebService/foundation/web"
	"log"
	"math/rand"
	"net/http"
)

type Check struct {
	log *log.Logger
}

func (c Check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%100 == 0 {
		return web.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
