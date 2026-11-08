package httpserver

import (
	"log/slog"
	"net/http"
	"strings"

	"server/internal/httpserver/handlers/data_index"
	"server/internal/httpserver/handlers/healthz"
	"server/internal/httpserver/handlers/not_found"
	"server/internal/httpserver/middleware"
	"server/internal/mod"
)

func NewHandler(log *slog.Logger, dd *mod.DataDir) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz.Handler)
	mux.HandleFunc("GET /{$}", data_index.New(dd, log))
	mux.HandleFunc("/", not_found.Handler)

	return middleware.Apply(mux,
		middleware.NewAccessLog(log, slog.LevelInfo, func() func(*http.Request) bool {
			// the following routes will NOT be logged
			skipMap := map[string]struct{}{
				"/healthz":              {},
				"/apple-touch-icon.png": {},
				"/favicon.ico":          {},
				"/robots.txt":           {},
			}

			return func(r *http.Request) bool {
				_, skip := skipMap[strings.ToLower(r.URL.Path)]

				return skip
			}
		}()),
	)
}
