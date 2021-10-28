package mid

import (
	"context"
	"github.com/DumanYessengali/ardanlabWebService/foundation/web"
	"log"
	"net/http"
	"time"
)

func Logger(log *log.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, ok := ctx.Value(web.KeyValue).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			log.Printf("%s : started    : %s %s -> %s",
				v.TraceID, r.Method, r.URL.Path, r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Printf("%s : completed  : %s %s -> %s (%d) (%s)",
				v.TraceID, r.Method, r.URL.Path, r.RemoteAddr, v.StatusCode, time.Since(v.Now))

			return err
		}

		return h
	}
	return m
}
