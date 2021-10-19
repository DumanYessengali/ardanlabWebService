package web

import (
	"context"
	"github.com/dimfeld/httptreemux"
	"net/http"
	"os"
	"syscall"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App ...
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

// NewApp ...
func NewApp(shutdown chan os.Signal) *App {
	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
	return &app
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle ...
func (a *App) Handle(method string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	a.ContextMux.Handle(method, path, h)
}
