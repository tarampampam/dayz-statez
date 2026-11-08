package middleware_test

import (
	"bufio"
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"server/internal/httpserver/middleware"
	"server/internal/testutil/assert"
)

// TestNewAccessLog is the primary table-driven test for the NewAccessLog middleware.
// Each case provides a pre-built *http.Request, an inner handler, an optional skipper,
// and a checkLog function that receives the captured log records for assertions.
func TestNewAccessLog(t *testing.T) {
	t.Parallel()

	t.Run("common cases", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			giveRequest        *http.Request
			giveLevel          slog.Level
			giveSkipper        func(*http.Request) bool
			giveHandler        http.HandlerFunc
			giveResponseWriter func() http.ResponseWriter // nil → httptest.NewRecorder()
			checkLog           func(t *testing.T, capture *logCapture)
		}{
			"2xx: logged at configured level": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/api/resource?foo=bar", nil),
				giveLevel:   slog.LevelInfo,
				giveHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) },
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
					rec := capture.records()[0]
					a := attrs(rec)
					assert.Equal(t, slog.LevelInfo, rec.Level)
					assert.Contains(t, rec.Message, "HTTP request processed")
					assert.Equal(t, int64(http.StatusOK), a["status_code"].Int64())
					assert.Equal(t, "/api/resource?foo=bar", a["uri"].String())
					_, hasDuration := a["duration"]
					assert.True(t, hasDuration)
				},
			},
			"3xx: logged at configured level": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/old", nil),
				giveLevel:   slog.LevelInfo,
				giveHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusMovedPermanently) },
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
					rec := capture.records()[0]
					assert.Equal(t, slog.LevelInfo, rec.Level)
					assert.Contains(t, rec.Message, "HTTP request processed")
					assert.Equal(t, int64(http.StatusMovedPermanently), attrs(rec)["status_code"].Int64())
				},
			},
			"4xx: logged at configured level": {
				giveRequest: httptest.NewRequest(http.MethodPost, "/protected", nil),
				giveLevel:   slog.LevelDebug,
				giveHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusUnauthorized) },
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
					rec := capture.records()[0]
					assert.Equal(t, slog.LevelDebug, rec.Level)
					assert.Contains(t, rec.Message, "HTTP request processed")
					assert.Equal(t, int64(http.StatusUnauthorized), attrs(rec)["status_code"].Int64())
				},
			},
			"5xx: logged at configured level": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/crash", nil),
				giveLevel:   slog.LevelDebug,
				giveHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusInternalServerError) },
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
					rec := capture.records()[0]
					assert.Equal(t, slog.LevelDebug, rec.Level)
					assert.Contains(t, rec.Message, "HTTP request processed")
					assert.Equal(t, int64(http.StatusInternalServerError), attrs(rec)["status_code"].Int64())
				},
			},
			"implicit 200 when body written without WriteHeader": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/implicit", nil),
				giveLevel:   slog.LevelInfo,
				giveHandler: func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("hello"))
				},
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
					a := attrs(capture.records()[0])
					assert.Equal(t, int64(http.StatusOK), a["status_code"].Int64())
					assert.Equal(t, int64(5), a["response_size"].Int64())
				},
			},
			"response_size accumulates across multiple Write calls": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/chunked", nil),
				giveLevel:   slog.LevelInfo,
				giveHandler: func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("foo"))
					_, _ = w.Write([]byte("bar"))
					_, _ = w.Write([]byte("baz"))
				},
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
					assert.Equal(t, int64(9), attrs(capture.records()[0])["response_size"].Int64())
				},
			},
			"all standard request fields are present in log entry": {
				giveRequest: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/api/v1/items", nil)
					req.Header.Set("User-Agent", "test-agent/2.0")
					req.Host = "api.example.com"

					return req
				}(),
				giveLevel:   slog.LevelInfo,
				giveHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusCreated) },
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
					a := attrs(capture.records()[0])
					assert.Equal(t, int64(http.StatusCreated), a["status_code"].Int64())
					assert.Equal(t, "POST", a["method"].String())
					assert.Equal(t, "test-agent/2.0", a["useragent"].String())
					assert.Equal(t, "api.example.com", a["host"].String())
					assert.True(t, a["remote_addr"].String() != "")
					_, hasDuration := a["duration"]
					assert.True(t, hasDuration)

					_, hasURI := a["uri"]
					assert.True(t, hasURI)
					assert.Equal(t, int64(0), a["response_size"].Int64())
				},
			},
			"skipper returns true: handler is called but nothing is logged": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/healthz", nil),
				giveLevel:   slog.LevelInfo,
				giveSkipper: func(_ *http.Request) bool { return true },
				giveHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) },
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 0, capture.len())
				},
			},
			"skipper returns false: request is logged normally": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/api", nil),
				giveLevel:   slog.LevelInfo,
				giveSkipper: func(_ *http.Request) bool { return false },
				giveHandler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) },
				checkLog: func(t *testing.T, capture *logCapture) {
					assert.Equal(t, 1, capture.len())
				},
			},
			"hijack is delegated to underlying ResponseWriter": {
				giveRequest: httptest.NewRequest(http.MethodGet, "/ws", nil),
				giveLevel:   slog.LevelInfo,
				giveResponseWriter: func() http.ResponseWriter {
					return &hijackableResponseWriter{ResponseRecorder: httptest.NewRecorder()}
				},
				giveHandler: func(w http.ResponseWriter, _ *http.Request) {
					hj, ok := w.(http.Hijacker)
					assert.True(t, ok)

					conn, _, err := hj.Hijack()
					assert.NoError(t, err)

					_ = conn.Close()
				},
				checkLog: func(t *testing.T, capture *logCapture) {
					// after Hijack the status is never written, so the middleware
					// logs with the default implicit 200.
					assert.Equal(t, 1, capture.len())
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				capture := &logCapture{}

				var handlerCalled bool

				inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handlerCalled = true

					if tc.giveHandler != nil {
						tc.giveHandler(w, r)
					}
				})

				handler := middleware.NewAccessLog(slog.New(capture), tc.giveLevel, tc.giveSkipper)(inner)

				w := http.ResponseWriter(httptest.NewRecorder())
				if tc.giveResponseWriter != nil {
					w = tc.giveResponseWriter()
				}

				assert.NotPanics(t, func() { handler.ServeHTTP(w, tc.giveRequest) })
				assert.True(t, handlerCalled)

				if tc.checkLog != nil {
					tc.checkLog(t, capture)
				}
			})
		}
	})

	t.Run("doubled invocation guard", func(t *testing.T) {
		t.Parallel()

		capture := &logCapture{}

		mw := middleware.NewAccessLog(slog.New(capture), slog.LevelInfo, nil)

		// wrap the inner handler twice with the same middleware constructor.
		handler := mw(mw(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		))

		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

		assert.Equal(t, 1, capture.len())
	})
}

// --------------------------------------------------------------------------------------------------------------------

// logCapture is a slog.Handler that records every emitted record for assertions. It accepts all levels.
type logCapture struct {
	mu   sync.Mutex
	recs []slog.Record
}

func (c *logCapture) Enabled(context.Context, slog.Level) bool { return true }
func (c *logCapture) WithAttrs([]slog.Attr) slog.Handler       { return c }
func (c *logCapture) WithGroup(string) slog.Handler            { return c }

func (c *logCapture) Handle(_ context.Context, r slog.Record) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.recs = append(c.recs, r.Clone())

	return nil
}

func (c *logCapture) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.recs)
}

func (c *logCapture) records() []slog.Record {
	c.mu.Lock()
	defer c.mu.Unlock()

	return append([]slog.Record(nil), c.recs...)
}

// attrs collects a record's attributes into a map keyed by attribute name.
func attrs(r slog.Record) map[string]slog.Value {
	m := make(map[string]slog.Value, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		m[a.Key] = a.Value

		return true
	})

	return m
}

// hijackableResponseWriter is a test double that adds http.Hijacker on top of httptest.ResponseRecorder so the
// Hijack-delegation path can be exercised.
type hijackableResponseWriter struct {
	*httptest.ResponseRecorder
}

func (h *hijackableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c1, c2 := net.Pipe()
	_ = c2.Close()

	return c1, bufio.NewReadWriter(bufio.NewReader(c1), bufio.NewWriter(c1)), nil
}
