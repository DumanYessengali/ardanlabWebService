package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Check struct {
	log *log.Logger
}

func (c Check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	c.log.Println(r, status)
	return json.NewEncoder(w).Encode(status)
}