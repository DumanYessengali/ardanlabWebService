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
		return errors.New("untrusted error")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	c.log.Println(r, status)
	return web.Respond(ctx, w, status, http.StatusOK)
}
