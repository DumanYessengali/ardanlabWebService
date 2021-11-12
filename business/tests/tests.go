package tests

import (
	"context"
	"github.com/DumanYessengali/ardanlabWebService/business/data/schema"
	"github.com/DumanYessengali/ardanlabWebService/foundation/database"
	"github.com/DumanYessengali/ardanlabWebService/foundation/web"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"testing"
	"time"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

var (
	dbImage = "postgres:13-alpine"
	dbPort  = "5432"
	dbArgs  = []string{"-e", "POSTGRES_PASSWORD=postgres"}
	AdminID = "7b5845ca-3436-444b-a128-4ec59481b23a"
	UserID  = "4a5794d8-1d87-4234-99ee-5e28ad6a4188"
)

func NewUint(t *testing.T) (*log.Logger, *sqlx.DB, func()) {
	c := startContainer(t, dbImage, dbPort, dbArgs...)

	cfg := database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	}

	db, err := database.Open(cfg)
	if err != nil {
		t.Fatalf("opening database connection %v", err)
	}

	t.Log("waiting for databse to be ready ...")

	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		dumpContainerLogs(t, c.ID)
		stopContainer(t, c.ID)
		t.Fatalf("database never ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c.ID)
		t.Fatalf("migrating error: %s", err)
	}

	teardown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c.ID)
	}

	log := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	return log, db, teardown
}

func Context() context.Context {
	values := web.Values{
		TraceID: uuid.New().String(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), web.KeyValues, &values)
}

func StringPointer(s string) *string {
	return &s
}

func IntPointer(i int) *int {
	return &i
}
